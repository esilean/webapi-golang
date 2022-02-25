package main

import (
	"backend/cmd/api/dtos"
	"backend/models"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
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

		movie = app.getPoster(movie)

		err := app.models.DB.AddMovie(movie)
		if err != nil {
			app.errorJSON(w, err)
			return
		}
	} else {

		movieExists, err := app.models.DB.GetMovieById(movie.Id)
		if err != nil {
			app.errorJSON(w, err)
			return
		}

		if movieExists != nil && movieExists.Poster == "" {
			movie = app.getPoster(movie)
		} else {
			movie.Poster = movieExists.Poster
		}

		err = app.models.DB.UpdateMovie(movie)
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

func (app *application) getPoster(movie models.Movie) models.Movie {

	type TheMovieDB struct {
		Page    int `json:"page"`
		Results []struct {
			Adult            bool    `json:"adult"`
			BackdropPath     string  `json:"backdrop_path"`
			GenreIds         []int   `json:"genre_ids"`
			ID               int     `json:"id"`
			OriginalLanguage string  `json:"original_language"`
			OriginalTitle    string  `json:"original_title"`
			Overview         string  `json:"overview"`
			Popularity       float64 `json:"popularity"`
			PosterPath       string  `json:"poster_path"`
			ReleaseDate      string  `json:"release_date"`
			Title            string  `json:"title"`
			Video            bool    `json:"video"`
			VoteAverage      float64 `json:"vote_average"`
			VoteCount        int     `json:"vote_count"`
		} `json:"results"`
		TotalPages   int `json:"total_pages"`
		TotalResults int `json:"total_results"`
	}

	client := &http.Client{}
	key := app.config.theMovieDB.key
	theUrl := "https://api.themoviedb.org/3/search/movie?api_key="
	theURI := theUrl + key + "&query=" + url.QueryEscape(movie.Title)

	req, err := http.NewRequest("GET", theURI, nil)
	if err != nil {
		app.logger.Println(err)
		return movie
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		app.logger.Println(err)
		return movie
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		app.logger.Println(err)
		return movie
	}

	var responseTheMovieDB TheMovieDB
	json.Unmarshal(bodyBytes, &responseTheMovieDB)

	if len(responseTheMovieDB.Results) > 0 {
		movie.Poster = responseTheMovieDB.Results[0].PosterPath
	}

	return movie
}
