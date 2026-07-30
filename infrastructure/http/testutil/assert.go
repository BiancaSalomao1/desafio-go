package testutil

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func AssertStatus(
	t *testing.T,
	rec *httptest.ResponseRecorder,
	expected int,
) {

	t.Helper()

	if rec.Code != expected {
		t.Fatalf(
			"expected status %d, got %d",
			expected,
			rec.Code,
		)
	}
}

func AssertContentTypeJSON(
	t *testing.T,
	rec *httptest.ResponseRecorder,
) {

	t.Helper()

	contentType := rec.Header().Get("Content-Type")

	if !strings.Contains(contentType, "application/json") {
		t.Fatalf(
			"expected application/json, got %s",
			contentType,
		)
	}
}

func AssertBodyContains(
	t *testing.T,
	rec *httptest.ResponseRecorder,
	expected string,
) {

	t.Helper()

	if !strings.Contains(
		rec.Body.String(),
		expected,
	) {
		t.Fatalf(
			"expected body to contain %q\nbody:\n%s",
			expected,
			rec.Body.String(),
		)
	}
}
