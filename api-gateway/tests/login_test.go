package tests

import (
	"net/http"
	"testing"
)

func TestLoginEndpoint(t *testing.T) {
	body := map[string]any{"email": "user@example.com", "password": "secret123"}
	w := doRequest(setupRouter(), http.MethodPost, "/api/login", body)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body=%s", w.Code, http.StatusOK, w.Body.String())
	}
}
