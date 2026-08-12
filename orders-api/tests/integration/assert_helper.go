/*
assert_helper

Responsabilidades:

- centralizar as principais asserções dos testes;
- reduzir repetição;
- melhorar legibilidade.
*/

package integration

import (
	"net/http"
	"testing"
)

func assertStatus(
	t *testing.T,
	resp *http.Response,
	expected int,
) {
	t.Helper()

	if resp.StatusCode != expected {
		t.Fatalf(
			"expected status %d got %d",
			expected,
			resp.StatusCode,
		)
	}
}

func assertCreated(
	t *testing.T,
	resp *http.Response,
) {
	t.Helper()
	assertStatus(t, resp, http.StatusCreated)
}

func assertOK(
	t *testing.T,
	resp *http.Response,
) {
	t.Helper()
	assertStatus(t, resp, http.StatusOK)
}

func assertNoContent(
	t *testing.T,
	resp *http.Response,
) {
	t.Helper()
	assertStatus(t, resp, http.StatusNoContent)
}

func assertBadRequest(
	t *testing.T,
	resp *http.Response,
) {
	t.Helper()
	assertStatus(t, resp, http.StatusBadRequest)
}

func assertUnauthorized(
	t *testing.T,
	resp *http.Response,
) {
	t.Helper()
	assertStatus(t, resp, http.StatusUnauthorized)
}

func assertNotFound(
	t *testing.T,
	resp *http.Response,
) {
	t.Helper()
	assertStatus(t, resp, http.StatusNotFound)
}

func assertConflict(t *testing.T, resp *http.Response) {
	t.Helper()

	if resp.StatusCode != http.StatusConflict {
		t.Fatalf(
			"expected status %d got %d",
			http.StatusConflict,
			resp.StatusCode,
		)
	}
}
