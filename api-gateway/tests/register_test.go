package tests

import (
	"net/http"
	"testing"
)

func TestRegisterEndpoint(t *testing.T) {
	body := map[string]any{"email": "user@example.com", "password": "secret123", "name": "John Doe"}
	w := doRequest(setupRouter(), http.MethodPost, "/api/register", body)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d, body=%s", w.Code, http.StatusCreated, w.Body.String())
	}
}
