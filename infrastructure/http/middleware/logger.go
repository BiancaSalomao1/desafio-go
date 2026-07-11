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
	"log"
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

func Logger(next http.Handler) http.Handler {

	return http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {

		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		start := time.Now()

		next.ServeHTTP(rw, r)

		log.Printf(
			"%s %s %d %s",
			r.Method,
			r.URL.Path,
			rw.statusCode,
			time.Since(start),
		)
	})
}
