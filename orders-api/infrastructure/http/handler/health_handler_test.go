package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	handler := NewHealthHandler()

	req := httptest.NewRequest(
		http.MethodGet,
		"/health",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.Health(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusOK,
			rec.Code,
		)
	}

	expected := `{"status":"UP"}` + "\n"

	if rec.Body.String() != expected {
		t.Errorf(
			"expected body %s, got %s",
			expected,
			rec.Body.String(),
		)
	}
}
