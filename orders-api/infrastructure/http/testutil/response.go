package testutil

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
)

func DecodeResponse[T any](
	t *testing.T,
	rec *httptest.ResponseRecorder,
) T {

	t.Helper()

	var response T

	if err := json.Unmarshal(
		rec.Body.Bytes(),
		&response,
	); err != nil {

		t.Fatalf(
			"cannot decode response: %v",
			err,
		)
	}

	return response
}
