package main

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) wrap(next http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		ctx := context.WithValue(r.Context(), httprouter.ParamsKey, ps)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func (app *application) routes() http.Handler {
	router := httprouter.New()
	secure := alice.New(app.checkToken)

	router.HandlerFunc(http.MethodGet, "/status", app.getApplicationStatus)

	router.HandlerFunc(http.MethodPost, "/v1/signin", app.Signin)

	router.HandlerFunc(http.MethodGet, "/v1/movie/:id", app.getMovieById)
	router.HandlerFunc(http.MethodGet, "/v1/movies", app.getMovies)
	router.HandlerFunc(http.MethodGet, "/v1/movies/genre/:id", app.getMoviesByGenreId)

	router.POST("/v1/admin/movie", app.wrap(secure.ThenFunc(app.upsertMovie)))
	router.DELETE("/v1/admin/movie/:id", app.wrap(secure.ThenFunc(app.deleteMovie)))
	//router.HandlerFunc(http.MethodDelete, "/v1/admin/movie/:id", app.deleteMovie)

	router.HandlerFunc(http.MethodGet, "/v1/genres", app.getGenres)

	return app.enableCORS(router)
}
