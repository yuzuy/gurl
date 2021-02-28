package main

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMakeDefaultHeader(t *testing.T) {
	cf := configFile{
		"127.0.0.1:8080": {
			Header: map[string]string{
				"Accept-Language": "en-US",
				"Content-Type":    "application/json",
			},
		},
		"127.0.0.1:8080/v1/*": {
			Header: map[string]string{
				"Authorization": "Basic foo",
				"Content-Type":  "x-www-form-urlencoded",
			},
		},
		"127.0.0.1:8080/v1/foo": {
			Header: map[string]string{
				"Accept-Charset": "utf-8",
				"Authorization":  "Basic bar",
			},
		},
		"127.0.0.1:8888": {
			Header: map[string]string{
				"Content-Type": "text/plain",
			},
		},
	}

	tests := []struct {
		url      string
		expected http.Header
	}{
		{
			url: "http://127.0.0.1:8080/v2/foo",
			expected: http.Header{
				"Accept-Language": {"en-US"},
				"Content-Type":    {"application/json"},
			},
		},
		{
			url: "http://127.0.0.1:8080/v1/bar",
			expected: http.Header{
				"Accept-Language": {"en-US"},
				"Authorization":   {"Basic foo"},
				"Content-Type":    {"x-www-form-urlencoded"},
			},
		},
		{
			url: "http://127.0.0.1:8080/v1/foo",
			expected: http.Header{
				"Accept-Charset":  {"utf-8"},
				"Accept-Language": {"en-US"},
				"Authorization":   {"Basic bar"},
				"Content-Type":    {"x-www-form-urlencoded"},
			},
		},
		{
			url: "http://127.0.0.1:8080/v1/foo/bar",
			expected: http.Header{
				"Accept-Language": {"en-US"},
				"Authorization":   {"Basic foo"},
				"Content-Type":    {"x-www-form-urlencoded"},
			},
		},
	}

	for _, tt := range tests {
		uri, err := url.Parse(tt.url)
		if err != nil {
			t.Fatalf("parsing tt.url failed. err=%s", err.Error())
		}
		got, err := makeDefaultHeader(uri, cf)
		if err != nil {
			t.Fatalf("makeDefaultHeader failed. url=%s, err=%s", tt.url, err)
		}

		if !cmp.Equal(got, tt.expected) {
			t.Errorf("makeDefaultHeader wrong. got=%v, expected=%v", got, tt.expected)
		}
	}
}
