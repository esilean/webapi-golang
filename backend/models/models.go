package models

import (
	"time"

	"gorm.io/gorm"
)

// Models is the wrapper for database
type Models struct {
	DB DBModel
}

// NewModels returns models with db pool
func NewModels(db *gorm.DB) Models {
	return Models{
		DB: DBModel{DB: db},
	}
}

type Movie struct {
	Id          int            `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Year        int            `json:"year"`
	ReleaseDate time.Time      `json:"release_date"`
	Runtime     int            `json:"runtime"`
	Rating      int            `json:"rating"`
	MPAARating  string         `json:"mpaa_rating"`
	Poster      string         `json:"poster"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	MovieGenre  map[int]string `json:"genres" gorm:"-"`
}

type Genre struct {
	Id        int       `json:"id"`
	Name      string    `json:"name" gorm:"column:genre_name"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type MovieGenre struct {
	Id        int       `json:"-"`
	MovieId   int       `json:"-"`
	GenreId   int       `json:"-"`
	Genre     Genre     `json:"genre"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type User struct {
	Id       int
	Email    string
	Password string
}
