package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newTestRouter() http.Handler {
	apiCfg := &apiConfig{}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/users", apiCfg.createUserHandler)
	return mux
}

func postJSON(t *testing.T, handler http.Handler, path string, body interface{}) *httptest.ResponseRecorder {
	t.Helper()
	b, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("failed to marshal body: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	return rr
}

func TestCreateUserHandler_MissingName(t *testing.T) {
	rr := postJSON(t, newTestRouter(), "/v1/users", map[string]string{})
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
	body := rr.Body.String()
	if !strings.Contains(body, "name is required") {
		t.Errorf("expected 'name is required' in body, got: %s", body)
	}
}

func TestCreateUserHandler_EmptyName(t *testing.T) {
	rr := postJSON(t, newTestRouter(), "/v1/users", map[string]string{"name": ""})
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestCreateUserHandler_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/v1/users", strings.NewReader("not json"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	newTestRouter().ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestLoggingMiddleware_CapuresStatus(t *testing.T) {
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	})
	wrapped := loggingMiddleware(inner)
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()
	wrapped.ServeHTTP(rr, req)
	if rr.Code != http.StatusTeapot {
		t.Errorf("expected 418, got %d", rr.Code)
	}
}

func TestReadinessHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/v1/healthz", nil)
	rr := httptest.NewRecorder()
	readinessHandler(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", ct)
	}
}
