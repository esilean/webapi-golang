package main

import (
	"backend/cmd/api/dtos"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (app *application) getMovieById(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.logger.Println(errors.New("invalid id parameter"))
		app.errorJSON(w, err)
		return
	}

	movie, err := app.models.DB.GetMovieById(id)
	if err != nil {
		app.logger.Println(err)
		app.errorJSON(w, err)
	}

	if err = app.writeJSON(w, http.StatusOK, movie, "movie"); err != nil {
		app.errorJSON(w, err)
	}
}

func (app *application) getMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := app.models.DB.GetAllMovies()
	if err != nil {
		app.logger.Println(err)
		app.errorJSON(w, err)
	}

	if err = app.writeJSON(w, http.StatusOK, movies, "movies"); err != nil {
		app.errorJSON(w, err)
	}
}

func (app *application) getMoviesByGenreId(w http.ResponseWriter, r *http.Request) {

	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.logger.Println(errors.New("invalid id parameter"))
		app.errorJSON(w, err)
		return
	}

	movies, err := app.models.DB.GetAllMovies(id)
	if err != nil {
		app.logger.Println(err)
		app.errorJSON(w, err)
	}

	if err = app.writeJSON(w, http.StatusOK, movies, "movies"); err != nil {
		app.errorJSON(w, err)
	}
}

func (app *application) getGenres(w http.ResponseWriter, r *http.Request) {
	genres, err := app.models.DB.GetAllGenres()
	if err != nil {
		app.logger.Println(err)
		app.errorJSON(w, err)
	}

	if err = app.writeJSON(w, http.StatusOK, genres, "genres"); err != nil {
		app.errorJSON(w, err)
	}
}

func (app *application) upsertMovie(w http.ResponseWriter, r *http.Request) {

	var movieRequest dtos.MovieRequest

	if err := json.NewDecoder(r.Body).Decode(&movieRequest); err != nil {
		app.errorJSON(w, err)
		return
	}

	movie := movieRequest.ToDomain()
	if movie.Id == 0 {
		err := app.models.DB.AddMovie(movie)
		if err != nil {
			app.errorJSON(w, err)
			return
		}
	} else {
		err := app.models.DB.UpdateMovie(movie)
		if err != nil {
			app.errorJSON(w, err)
			return
		}
	}

	if err := app.writeNoContent(w, http.StatusOK); err != nil {
		app.errorJSON(w, err)
	}
}

func (app *application) deleteMovie(w http.ResponseWriter, r *http.Request) {

	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.logger.Println(errors.New("invalid id parameter"))
		app.errorJSON(w, err)
		return
	}

	err = app.models.DB.DeleteMovie(id)
	if err != nil {
		app.logger.Println(err)
		app.errorJSON(w, err)
	}

	if err := app.writeNoContent(w, http.StatusNoContent); err != nil {
		app.errorJSON(w, err)
	}
}
