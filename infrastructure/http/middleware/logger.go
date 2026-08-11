package middleware

/*
Função Logger

Responsabilidades:

- registrar início da requisição HTTP;
- registrar conclusão da requisição;
- registrar serviço;
- registrar operação;
- registrar resultado;
- registrar status HTTP;
- registrar tempo de execução;
- registrar erro quando houver.

Métodos:

- Logger()
*/

import (
	"log/slog"
	"net/http"
	"time"
)

const serviceName = "orders-api"

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(statusCode int) {

	rw.statusCode = statusCode

	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *responseWriter) Write(body []byte) (int, error) {

	if rw.statusCode == 0 {
		rw.statusCode = http.StatusOK
	}

	return rw.ResponseWriter.Write(body)
}

func Logger(next http.Handler) http.Handler {

	return http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {

		start := time.Now()

		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		operation := r.Method + " " + r.URL.Path

		slog.Info(
			"request started",
			"service", serviceName,
			"operation", operation,
		)

		next.ServeHTTP(rw, r)

		duration := time.Since(start)

		result := "success"

		if rw.statusCode >= http.StatusBadRequest {
			result = "error"
		}

		attrs := []any{
			"service", serviceName,
			"operation", operation,
			"result", result,
			"status", rw.statusCode,
			"duration", duration.String(),
		}

		if rw.statusCode >= http.StatusInternalServerError {
			attrs = append(
				attrs,
				"error", http.StatusText(rw.statusCode),
			)
		}

		slog.Info(
			"request completed",
			attrs...,
		)
	})
}
