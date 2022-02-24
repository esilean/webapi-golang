package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/status", app.getApplicationStatus)

	router.HandlerFunc(http.MethodGet, "/v1/movie/:id", app.getMovieById)
	router.HandlerFunc(http.MethodGet, "/v1/movies", app.getMovies)
	router.HandlerFunc(http.MethodGet, "/v1/movies/genre/:id", app.getMoviesByGenreId)

	router.HandlerFunc(http.MethodPost, "/v1/admin/movie", app.upsertMovie)
	router.HandlerFunc(http.MethodDelete, "/v1/admin/movie/:id", app.deleteMovie)

	router.HandlerFunc(http.MethodGet, "/v1/genres", app.getGenres)

	return app.enableCORS(router)
}
