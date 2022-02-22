package models

import (
	"gorm.io/gorm"
)

type DBModel struct {
	DB *gorm.DB
}

func (m *DBModel) GetById(id int) (*Movie, error) {

	var movie Movie

	m.DB.First(&movie, id)

	genres, err := m.getGenresByMovieId(id)
	if err != nil {
		return nil, err
	}

	movie.MovieGenre = genres

	return &movie, nil
}

func (m *DBModel) GetAll() ([]*Movie, error) {

	var movies []*Movie

	m.DB.Find(&movies)

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
