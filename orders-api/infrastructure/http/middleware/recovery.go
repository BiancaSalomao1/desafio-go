package middleware

/*
Função Recovery

Responsabilidades:
- recuperar panics da aplicação;
- registrar o erro no console;
- retornar erro HTTP 500.

Métodos:
- Recovery()
*/

import (
	"fmt"
	"net/http"

	"desafio-go/orders-api/infrastructure/http/httpx"
)

func Recovery(next http.Handler) http.Handler {

	return http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {

		defer func() {

			if err := recover(); err != nil {

				fmt.Printf("PANIC: %v\n", err)

				httpx.WriteJSON(
					w,
					http.StatusInternalServerError,
					map[string]string{
						"error": "internal server error",
					},
				)
			}

		}()

		next.ServeHTTP(w, r)
	})
}
