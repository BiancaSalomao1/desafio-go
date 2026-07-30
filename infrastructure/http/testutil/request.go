package testutil

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

func NewJSONRequest(
	method string,
	url string,
	body any,
) (*http.Request, *httptest.ResponseRecorder) {

	var buffer bytes.Buffer

	if body != nil {
		_ = json.NewEncoder(&buffer).Encode(body)
	}

	req := httptest.NewRequest(
		method,
		url,
		&buffer,
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	rec := httptest.NewRecorder()

	return req, rec
}
