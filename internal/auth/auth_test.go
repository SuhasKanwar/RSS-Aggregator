package auth

import (
	"net/http"
	"testing"
)

func TestGetAPIKey(t *testing.T) {
	tests := []struct {
		name        string
		headerValue string
		wantKey     string
		wantErr     bool
	}{
		{
			name:        "valid api key",
			headerValue: "ApiKey abc123",
			wantKey:     "abc123",
			wantErr:     false,
		},
		{
			name:        "missing header",
			headerValue: "",
			wantKey:     "",
			wantErr:     true,
		},
		{
			name:        "wrong scheme",
			headerValue: "Bearer abc123",
			wantKey:     "",
			wantErr:     true,
		},
		{
			name:        "no space separator",
			headerValue: "ApiKeyabc123",
			wantKey:     "",
			wantErr:     true,
		},
		{
			name:        "too many parts",
			headerValue: "ApiKey abc 123",
			wantKey:     "",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headers := http.Header{}
			if tt.headerValue != "" {
				headers.Set("Authorization", tt.headerValue)
			}
			got, err := GetAPIKey(headers)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAPIKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.wantKey {
				t.Errorf("GetAPIKey() = %v, want %v", got, tt.wantKey)
			}
		})
	}
}
