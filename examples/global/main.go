package main

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/maltsev25/shutdown"
)

func main() {
	mux := http.NewServeMux()

	httpSrv := http.Server{
		Addr:    ":8888",
		Handler: mux,
	}

	shutdown.InitGlobal()

	shutdown.MustAdd("some_service", func(ctx context.Context) {
		log.Println("some_service shutdown success")
	})

	shutdown.MustAdd("other_service", func(ctx context.Context) {
		log.Println("other_service shutdown success")
	}, "some_service")

	shutdown.MustAdd("http_server", func(ctx context.Context) {
		if err := httpSrv.Shutdown(ctx); err != nil {
			log.Println("failed to shut down http_server")

			return
		}

		log.Println("http_server shutdown success")
	}, "other_service")

	if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}

	if err := shutdown.Wait(); err != nil {
		panic(err)
	}
}
