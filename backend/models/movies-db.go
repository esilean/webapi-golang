package models

import (
	"context"
	"database/sql"
	"time"
)

type DBModel struct {
	DB *sql.DB
}

func (m *DBModel) GetById(id int) (*Movie, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	queryMovie := `select id, title, description, year, release_date, rating, runtime, mpaa_rating, created_at, updated_at
			  from movies
			  where id = $1
			 `

	row := m.DB.QueryRowContext(ctx, queryMovie, id)

	var movie Movie

	err := row.Scan(
		&movie.Id,
		&movie.Title,
		&movie.Description,
		&movie.Year,
		&movie.ReleaseDate,
		&movie.Rating,
		&movie.Runtime,
		&movie.MPAARating,
		&movie.CreatedAt,
		&movie.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	queryGenres := `select 
				a.id, a.movie_id, a.genre_id, b.genre_name name, a.created_at, a.updated_at
			from movies_genres a
			left join genres b on (a.genre_id = b.id)
			where 
				a.movie_id = $1
			`

	rows, err := m.DB.QueryContext(ctx, queryGenres, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	genres := make(map[int]string, 0)
	for rows.Next() {
		var mg MovieGenre

		err := rows.Scan(
			&mg.Id,
			&mg.MovieId,
			&mg.GenreId,
			&mg.Genre.Name,
			&mg.CreatedAt,
			&mg.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		genres[mg.Id] = mg.Genre.Name
	}

	movie.MovieGenre = genres

	return &movie, nil
}

func (m *DBModel) GetAll() ([]*Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	queryMovies := `select 
				id, title, description, year, release_date, rating, runtime, mpaa_rating, created_at, updated_at 
				from movies order by title
			 `

	rows, err := m.DB.QueryContext(ctx, queryMovies)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	movies := make([]*Movie, 0)
	for rows.Next() {
		var movie Movie

		err := rows.Scan(
			&movie.Id,
			&movie.Title,
			&movie.Description,
			&movie.Year,
			&movie.ReleaseDate,
			&movie.Rating,
			&movie.Runtime,
			&movie.MPAARating,
			&movie.CreatedAt,
			&movie.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		queryGenres := `select 
					a.id, a.movie_id, a.genre_id, b.genre_name name, a.created_at, a.updated_at
				from movies_genres a
				left join genres b on (a.genre_id = b.id)
				where 
					a.movie_id = $1
				`

		rowsGenre, err := m.DB.QueryContext(ctx, queryGenres, &movie.Id)
		if err != nil {
			return nil, err
		}

		genres := make(map[int]string, 0)
		for rowsGenre.Next() {
			var mg MovieGenre

			err := rowsGenre.Scan(
				&mg.Id,
				&mg.MovieId,
				&mg.GenreId,
				&mg.Genre.Name,
				&mg.CreatedAt,
				&mg.UpdatedAt,
			)

			if err != nil {
				return nil, err
			}

			genres[mg.Id] = mg.Genre.Name
		}
		rowsGenre.Close()

		movie.MovieGenre = genres

		movies = append(movies, &movie)
	}

	return movies, nil
}
