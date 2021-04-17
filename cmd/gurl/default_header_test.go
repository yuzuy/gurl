package main

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPattern_match(t *testing.T) {
	tests := []struct {
		name   string
		p      pattern
		target string
		want   bool
	}{
		{
			name:   "match completely",
			p:      "localhost:8080",
			target: "http://localhost:8080",
			want:   true,
		},
		{
			name:   "match using a wildcard",
			p:      "localhost:8080/*",
			target: "http://localhost:8080/foo/bar",
			want:   true,
		},
		{
			name:   "match using a wildcard as a path variable",
			p:      "localhost:8080/v1/users/*/repositories",
			target: "http://localhost:8080/v1/users/foo/repositories",
			want:   true,
		},
		{
			name:   "match using multiple wildcards",
			p:      "localhost:8080/v1/users/*/repositories/*/issues",
			target: "http://localhost:8080/v1/users/foo/repositories/gurl/issues",
			want:   true,
		},
		{
			name:   "match with the pattern that is only wildcard",
			p:      "*",
			target: "localhost:8080",
			want:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uri, err := url.Parse(tt.target)
			if err != nil {
				t.Errorf("failed to parse tt.target(=%s): %+v", tt.target, err)
				return
			}
			m, err := tt.p.match(uri)
			if err != nil {
				t.Errorf("failed to tt.p.match with tt.target(=%s): %+v", tt.target, err)
				return
			}
			if m != tt.want {
				t.Errorf("pattern.match wrong. target=%s, want=%t, got=%t", tt.target, tt.want, m)
			}
		})
	}
}

func TestDefaultHeaderList_set(t *testing.T) {
	dhl := defaultHeaderList{
		"localhost:8080": {
			"Accept-Language": "en-US",
			"Content-Type":    "application/json",
		},
		"localhost:8080/v1/*": {
			"Authorization": "Basic foo",
			"Content-Type":  "x-www-form-urlencoded",
		},
		"localhost:8080/v1/bar": {
			"Accept-Charset": "utf-8",
			"Authorization":  "Basic bar",
		},
		"localhost:8888": {
			"Content-Type": "text/plain",
		},
	}

	tests := []struct {
		name     string
		url      string
		expected http.Header
	}{
		{
			name: "set header",
			url:  "http://localhost:8080/v2/foo",
			expected: http.Header{
				"Accept-Language": {"en-US"},
				"Content-Type":    {"application/json"},
			},
		},
		{
			name: "the default header for the deeper path has priority",
			url:  "http://localhost:8080/v1/foo",
			expected: http.Header{
				"Accept-Language": {"en-US"},
				"Authorization":   {"Basic foo"},
				"Content-Type":    {"x-www-form-urlencoded"},
			},
		},
		{
			name: "the default header for the more detailed path has priority",
			url:  "http://localhost:8080/v1/bar",
			expected: http.Header{
				"Accept-Charset":  {"utf-8"},
				"Accept-Language": {"en-US"},
				"Authorization":   {"Basic bar"},
				"Content-Type":    {"x-www-form-urlencoded"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, tt.url, nil)
			if err != nil {
				t.Errorf("failed to call http.NewRequest. args=(\"GET\", %s, nil)", tt.url)
				return
			}
			err = dhl.set(req)
			if err != nil {
				t.Errorf("makeHeaderFromDefaultHeader failed. url=%s, err=%s", tt.url, err)
				return
			}

			if !cmp.Equal(req.Header, tt.expected) {
				t.Errorf("defaultHeaderList.set wrong. got=%v, expected=%v", req.Header, tt.expected)
			}
		})
	}
}
