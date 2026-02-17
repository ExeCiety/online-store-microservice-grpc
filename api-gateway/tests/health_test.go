package tests

import (
	"net/http"
	"testing"
)

func TestHealthEndpoint(t *testing.T) {
	w := doRequest(setupRouter(), http.MethodGet, "/health", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
}
