/*
response_helper

Responsabilidades:

- executar requisições HTTP;
- fechar response.Body automaticamente.
*/

package integration

import (
	"net/http"
	"testing"
)

func doRequest(
	t *testing.T,
	req *http.Request,
) *http.Response {

	t.Helper()

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Fatal(err)
	}

	return resp
}
