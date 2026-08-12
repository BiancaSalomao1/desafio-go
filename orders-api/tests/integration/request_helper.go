package integration

/*
request_helper

Responsabilidades:

- criar requisições HTTP autenticadas;
- adicionar automaticamente o JWT.

Este arquivo evita repetição nos testes.
*/

import (
	"io"
	"net/http"
)

func authenticatedRequest(
	method string,
	url string,
	token string,
	body io.Reader,
) (*http.Request, error) {

	req, err := http.NewRequest(
		method,
		url,
		body,
	)

	if err != nil {
		return nil, err
	}

	req.Header.Set(
		"Authorization",
		"Bearer "+token,
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	return req, nil
}
