package testutil

import (
	"net/http"
	"net/http/httptest"
)

type RequestBuilder struct {
	method string
	url    string
	body   any
}

func NewRequestBuilder() *RequestBuilder {

	return &RequestBuilder{}
}

func (b *RequestBuilder) Method(
	method string,
) *RequestBuilder {

	b.method = method
	return b
}

func (b *RequestBuilder) URL(
	url string,
) *RequestBuilder {

	b.url = url
	return b
}

func (b *RequestBuilder) Body(
	body any,
) *RequestBuilder {

	b.body = body
	return b
}

func (b *RequestBuilder) Build() (*http.Request, *httptest.ResponseRecorder) {
	return NewJSONRequest(
		b.method,
		b.url,
		b.body,
	)
}
