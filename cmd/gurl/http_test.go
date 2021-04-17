package main

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMakeHeaderFromDefaultHeader(t *testing.T) {
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
			uri, err := url.Parse(tt.url)
			if err != nil {
				t.Errorf("parsing tt.url failed. err=%s", err.Error())
				return
			}
			got, err := makeHeaderFromDefaultHeader(uri, dhl)
			if err != nil {
				t.Errorf("makeHeaderFromDefaultHeader failed. url=%s, err=%s", tt.url, err)
				return
			}

			if !cmp.Equal(got, tt.expected) {
				t.Errorf("makeHeaderFromDefaultHeader wrong. got=%v, expected=%v", got, tt.expected)
			}
		})
	}
}

func TestParseHeader(t *testing.T) {
	tests := []struct {
		name        string
		h           string
		expectedKey string
		expectedVal string
		wantErr     bool
	}{
		{
			name:        "parse valid header",
			h:           "Content-Type: application/json",
			expectedKey: "Content-Type",
			expectedVal: "application/json",
			wantErr:     false,
		},
		{
			name:    "parse invalid header",
			h:       "Content-Type application/json",
			wantErr: true,
		},
		{
			name:        "allow colon in value",
			h:           "Content-Type: foo-format:v1",
			expectedKey: "Content-Type",
			expectedVal: "foo-format:v1",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, val, err := parseHeader(tt.h)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseHeader wrong. header=%s, wantErr=%t, gotErr=%v)", tt.h, tt.wantErr, err)
				return
			}
			if key != tt.expectedKey {
				t.Errorf("parseHeader wrong. expectedKey=%s, gotKey=%s", tt.expectedKey, key)
				return
			}
			if val != tt.expectedVal {
				t.Errorf("parseHeader wrong. expectedVal=%s, gotVal=%s", tt.expectedVal, val)
			}
		})
	}
}
