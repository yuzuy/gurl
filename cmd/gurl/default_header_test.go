package main

import (
	"net/url"
	"testing"
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
			}
			m, err := tt.p.match(uri)
			if err != nil {
				t.Errorf("failed to tt.p.match with tt.target(=%s): %+v", tt.target, err)
			}
			if m != tt.want {
				t.Errorf("pattern.match wrong. target=%s, want=%t, got=%t", tt.target, tt.want, m)
			}
		})
	}
}
