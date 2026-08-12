package httpx

/*
Funções auxiliares para respostas HTTP.

Responsabilidades:
- retornar respostas JSON;
- retornar erros em formato JSON;
- retornar mensagens de sucesso.

Funções:
- WriteJSON()
- WriteError()
*/

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func WriteJSON(
	w http.ResponseWriter,
	status int,
	data any,
) {

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(data)
}

func WriteError(
	w http.ResponseWriter,
	status int,
	err error,
) {

	log.Printf("WRITE ERROR: %T - %q", err, err.Error())

	WriteJSON(
		w,
		status,
		ErrorResponse{
			Error: errorMessage(err),
		},
	)
}

func WriteStatus(
	w http.ResponseWriter,
	status int,
) {
	w.WriteHeader(status)
}
