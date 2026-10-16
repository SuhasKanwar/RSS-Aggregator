package main

import (
	"testing"
	"time"
)

func TestParsePubDate(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "RFC1123Z format",
			input:   "Thu, 16 Oct 2026 10:00:00 +0530",
			wantErr: false,
		},
		{
			name:    "RFC1123 format",
			input:   "Thu, 16 Oct 2026 10:00:00 IST",
			wantErr: false,
		},
		{
			name:    "RFC822Z format",
			input:   "16 Oct 26 10:00 +0530",
			wantErr: false,
		},
		{
			name:    "RFC822 format",
			input:   "16 Oct 26 10:00 IST",
			wantErr: false,
		},
		{
			name:    "RFC3339 format",
			input:   "2026-10-16T10:00:00Z",
			wantErr: false,
		},
		{
			name:    "custom ISO with Z",
			input:   "2026-10-16T10:00:05Z",
			wantErr: false,
		},
		{
			name:    "leading/trailing whitespace",
			input:   "  Thu, 16 Oct 2026 10:00:00 +0530  ",
			wantErr: false,
		},
		{
			name:    "unrecognised format",
			input:   "October 16 2026",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parsePubDate(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parsePubDate(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if err == nil && got.IsZero() {
				t.Errorf("parsePubDate(%q) returned zero time without error", tt.input)
			}
			if err == nil && got.Location() != time.UTC {
				t.Errorf("parsePubDate(%q) = %v, want UTC", tt.input, got.Location())
			}
		})
	}
}
