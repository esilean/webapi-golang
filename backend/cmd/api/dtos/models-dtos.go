package dtos

import (
	"backend/models"
	"strconv"
	"time"
)

type MovieRequest struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Year        string `json:"year"`
	ReleaseDate string `json:"release_date"`
	Runtime     string `json:"runtime"`
	Rating      string `json:"rating"`
	MPAARating  string `json:"mpaa_rating"`
}

func (r *MovieRequest) ToDomain() models.Movie {

	var movie models.Movie
	movie.Id, _ = strconv.Atoi(r.Id)
	movie.Title = r.Title
	movie.Description = r.Description
	movie.ReleaseDate, _ = time.Parse("2006-01-02", r.ReleaseDate)
	movie.Year = movie.ReleaseDate.Year()
	movie.Runtime, _ = strconv.Atoi(r.Runtime)
	movie.Rating, _ = strconv.Atoi(r.Rating)
	movie.MPAARating = r.MPAARating
	movie.CreatedAt = time.Now()
	movie.UpdatedAt = time.Now()

	return movie
}
