package main

import "testing"

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
