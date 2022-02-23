package models

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type DBModel struct {
	DB *sql.DB
}

func (m *DBModel) GetMovieById(id int) (*Movie, error) {

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

	genres, err := m.getMovieGenres(id, ctx)
	if err != nil {
		return nil, err
	}

	movie.MovieGenre = genres

	return &movie, nil
}

func (m *DBModel) GetAllMovies(genreIds ...int) ([]*Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	where := ""
	if len(genreIds) > 0 {
		where = fmt.Sprintf("where id in (select movie_id from movies_genres where genre_id = %d)", genreIds[0])
	}

	queryMovies := fmt.Sprintf(`select 
								id, title, description, year, release_date, rating, runtime, 
								mpaa_rating, created_at, updated_at 
								from movies 
								%s
								order by title`, where)

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

		genres, err := m.getMovieGenres(movie.Id, ctx)
		if err != nil {
			return nil, err
		}

		movie.MovieGenre = genres

		movies = append(movies, &movie)
	}

	return movies, nil
}

func (m *DBModel) getMovieGenres(id int, ctx context.Context) (map[int]string, error) {

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

	return genres, nil
}

func (m *DBModel) GetAllGenres() ([]*Genre, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	queryGenres := `select a.id, a.genre_name name, a.created_at, a.updated_at
					from genres a
					order by a.genre_name
 					`

	rows, err := m.DB.QueryContext(ctx, queryGenres)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	genres := make([]*Genre, 0)
	for rows.Next() {

		var genre Genre
		err := rows.Scan(
			&genre.Id,
			&genre.Name,
			&genre.CreatedAt,
			&genre.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		genres = append(genres, &genre)
	}

	return genres, nil

}

func (m *DBModel) AddMovie(movie Movie) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stmt := `insert into movies 
			 (title, description, year, release_date, runtime, rating, mpaa_rating, created_at, updated_at)
			 values 
			 ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := m.DB.ExecContext(ctx, stmt,
		movie.Title,
		movie.Description,
		movie.Year,
		movie.ReleaseDate,
		movie.Runtime,
		movie.Rating,
		movie.MPAARating,
		movie.CreatedAt,
		movie.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (m *DBModel) UpdateMovie(movie Movie) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stmt := `update movies
			 set
				title = $1,
				description = $2,
				year = $3,
				release_date = $4,
				runtime = $5,
				rating = $6,
				mpaa_rating = $7,
				updated_at = $8
			 where
				id = $9
			 `

	_, err := m.DB.ExecContext(ctx, stmt,
		movie.Title,
		movie.Description,
		movie.Year,
		movie.ReleaseDate,
		movie.Runtime,
		movie.Rating,
		movie.MPAARating,
		movie.UpdatedAt,
		movie.Id,
	)
	if err != nil {
		return err
	}

	return nil
}
