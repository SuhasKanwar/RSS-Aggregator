package auth

import (
	"errors"
	"net/http"
	"strings"
)

// GetAPIKey extracts an API Key from
// the headers of an HTTP request
// Example:
// Authorization: ApiKey { api_key }
func GetAPIKey(headers http.Header) (string, error) {
	val := headers.Get("authorization")
	if val == "" {
		return "", errors.New("authorization header is missing")
	}

	vals := strings.Split(val, " ")
	if len(vals) != 2 || vals[0] != "ApiKey" {
		return "", errors.New("authorization header is malformed")
	}

	return vals[1], nil
}
