package middleware

/*
Função Logger

Responsabilidades:
- registrar requisições HTTP;
- registrar status HTTP;
- registrar tempo de execução.

Métodos:
- Logger()
*/

import (
	"log/slog"
	"net/http"
	"time"
)

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

		slog.Info(
			"request started",
			"method", r.Method,
			"path", r.URL.Path,
		)

		next.ServeHTTP(rw, r)

		duration := time.Since(start)

		slog.Info(
			"request completed",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rw.statusCode,
			"duration", duration.String(),
		)
	})
}
