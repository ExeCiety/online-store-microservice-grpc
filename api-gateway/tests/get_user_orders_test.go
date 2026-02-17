package tests

import (
	"net/http"
	"testing"
)

func TestGetOrdersByUserIDEndpoint(t *testing.T) {
	w := doRequest(setupRouter(), http.MethodGet, "/api/users/4e427d78-58c5-4f78-bfc1-e2c196e0b506/orders", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body=%s", w.Code, http.StatusOK, w.Body.String())
	}
}
