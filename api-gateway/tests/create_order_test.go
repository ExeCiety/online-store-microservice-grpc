package tests

import (
	"net/http"
	"testing"
)

func TestCreateOrderEndpoint(t *testing.T) {
	body := map[string]any{"user_id": "4e427d78-58c5-4f78-bfc1-e2c196e0b506", "product_name": "Laptop", "quantity": 1, "total_price": 15000000}
	w := doRequest(setupRouter(), http.MethodPost, "/api/orders", body)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d, body=%s", w.Code, http.StatusCreated, w.Body.String())
	}
}
