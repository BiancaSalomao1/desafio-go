package httpx

/*
Funções auxiliares para leitura de requisições HTTP.

Responsabilidades:
- ler JSON da requisição;
- validar JSON.

Funções:
- ReadJSON()
*/

import (
	"encoding/json"
	"errors"
	"net/http"
)

func ReadJSON(
	r *http.Request,
	dst any,
) error {

	if r.Body == nil {
		return errors.New("request body is empty")
	}

	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)

	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		return err
	}

	// Garante que exista apenas um objeto JSON
	if decoder.More() {
		return errors.New("request body must contain a single JSON object")
	}

	return nil
}
