package main

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMakeHeaderFromDefaultHeader(t *testing.T) {
	dhs := defaultHeaders{
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
			uri, err := url.Parse(tt.url)
			if err != nil {
				t.Errorf("parsing tt.url failed. err=%s", err.Error())
			}
			got, err := makeHeaderFromDefaultHeader(uri, dhs)
			if err != nil {
				t.Errorf("makeHeaderFromDefaultHeader failed. url=%s, err=%s", tt.url, err)
			}

			if !cmp.Equal(got, tt.expected) {
				t.Errorf("makeHeaderFromDefaultHeader wrong. got=%v, expected=%v", got, tt.expected)
			}
		})
	}
}
