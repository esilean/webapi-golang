package models

import (
	"gorm.io/gorm"
)

type DBModel struct {
	DB *gorm.DB
}

func (m *DBModel) GetMovieById(id int) (*Movie, error) {

	var movie Movie

	m.DB.First(&movie, id)

	genres, err := m.getGenresByMovieId(id)
	if err != nil {
		return nil, err
	}

	movie.MovieGenre = genres

	return &movie, nil
}

func (m *DBModel) GetAllMovies(genreIds ...int) ([]*Movie, error) {

	var movies []*Movie

	if len(genreIds) > 0 {
		subQuery := m.DB.Select("movie_id").Where("genre_id = ?", genreIds[0]).Table("movies_genres")
		m.DB.Where("id in (?)", subQuery).Find(&movies)
	} else {
		m.DB.Find(&movies)
	}

	for _, movie := range movies {

		genres, err := m.getGenresByMovieId(movie.Id)
		if err != nil {
			return nil, err
		}

		movie.MovieGenre = genres
	}

	return movies, nil
}

func (m *DBModel) getGenresByMovieId(movieId int) (map[int]string, error) {

	genres := make(map[int]string, 0)
	rows, err := m.DB.Table("movies_genres a").
		Select("a.genre_id id, b.genre_name name").
		Joins("left join genres b on (a.genre_id = b.id)").
		Where("a.movie_id = ?", movieId).
		Rows()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var mg MovieGenre

		err := rows.Scan(
			&mg.Id,
			&mg.Genre.Name,
		)

		if err != nil {
			return nil, err
		}

		genres[mg.Id] = mg.Genre.Name
	}

	return genres, nil
}

func (m *DBModel) GetAllGenres() ([]*Genre, error) {

	var genres []*Genre

	m.DB.Find(&genres)

	return genres, nil
}

func (m *DBModel) AddMovie(movie Movie) error {

	result := m.DB.Create(&movie)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (m *DBModel) UpdateMovie(movie Movie) error {

	result := m.DB.Save(&movie)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (m *DBModel) DeleteMovie(id int) error {

	result := m.DB.Delete(&Movie{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
