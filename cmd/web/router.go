package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) router() http.Handler {
	//TODO:Implement notFound, methodNotAllowed, internal server responses
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/", app.home)

	return router
}

//test

//test2 buzuk
