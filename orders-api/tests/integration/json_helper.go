/*
json_helper

Responsabilidades:

- serializar estruturas;
- desserializar respostas;
*/

package integration

import (
	"encoding/json"
	"io"
	"testing"
)

func mustMarshal(
	t *testing.T,
	value any,
) []byte {

	t.Helper()

	data, err := json.Marshal(value)

	if err != nil {
		t.Fatal(err)
	}

	return data
}

func decodeResponse[T any](
	t *testing.T,
	reader io.Reader,
) T {

	t.Helper()

	var response T

	if err := json.NewDecoder(reader).Decode(&response); err != nil {
		t.Fatal(err)
	}

	return response
}
