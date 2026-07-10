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

	WriteJSON(
		w,
		status,
		ErrorResponse{
			Error: err.Error(),
		},
	)
}

/*
WriteStatus

Responsabilidades:
- retornar apenas o status HTTP.

Funções:
- WriteStatus()
*/

func WriteStatus(
	w http.ResponseWriter,
	status int,
) {
	w.WriteHeader(status)
}
