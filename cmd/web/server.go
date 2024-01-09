package main

import (
	"fmt"
	"net/http"
)

func (app *application) serve() error {
	srv := http.Server{
		Addr:    fmt.Sprintf("localhost:%d", app.cfg.port),
		Handler: app.router(),
	}
	fmt.Printf("server listens on port %d\n", app.cfg.port)
	err := srv.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}
