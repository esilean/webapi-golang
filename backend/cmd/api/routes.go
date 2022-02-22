package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/status", app.getApplicationStatus)

	router.HandlerFunc(http.MethodGet, "/v1/movie/:id", app.getMovieById)
	router.HandlerFunc(http.MethodGet, "/v1/movies", app.getMovies)

	return router
}
