package app

import (
	"log/slog"
	"net/http"
	"net/http/pprof"
)

func StartPprofServer(addr string) *http.Server {

	mux := http.NewServeMux()

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	mux.Handle(
		"/debug/pprof/goroutine",
		pprof.Handler("goroutine"),
	)

	mux.Handle(
		"/debug/pprof/heap",
		pprof.Handler("heap"),
	)

	mux.Handle(
		"/debug/pprof/allocs",
		pprof.Handler("allocs"),
	)

	mux.Handle(
		"/debug/pprof/block",
		pprof.Handler("block"),
	)

	mux.Handle(
		"/debug/pprof/mutex",
		pprof.Handler("mutex"),
	)

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {

		slog.Info(
			"pprof server started",
			"address", addr,
		)

		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {

			slog.Error(
				"pprof server stopped",
				"error", err,
			)
		}
	}()

	return server
}
