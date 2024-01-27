package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) router() http.Handler {
	//TODO:Implement internal server responses
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	fileServer := http.FileServer(http.Dir("../../ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodPost, "/user/register", app.createUserHandlerPost)
	router.HandlerFunc(http.MethodGet, "/user/register", app.createUserHandlerGet)

	router.HandlerFunc(http.MethodGet, "/post/create", app.createPostHandlerGet)
	router.HandlerFunc(http.MethodPost, "/post/create", app.createPostHandlerPost)
	router.HandlerFunc(http.MethodGet, "/post/view/:id", app.showPostHandler)
	router.HandlerFunc(http.MethodPatch, "/post/view/:id", app.updatePostHandler)
	router.HandlerFunc(http.MethodDelete, "/post/view/:id", app.deletePostHandler)
	return router
}
