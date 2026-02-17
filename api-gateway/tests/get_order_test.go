package tests

import (
	"net/http"
	"testing"
)

func TestGetOrderByIDEndpoint(t *testing.T) {
	w := doRequest(setupRouter(), http.MethodGet, "/api/orders/8f328abb-4ae4-493b-a460-a63f1206b2f3", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body=%s", w.Code, http.StatusOK, w.Body.String())
	}
}
