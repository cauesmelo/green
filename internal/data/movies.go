package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/cauesmelo/green/internal/validator"
	"github.com/lib/pq"
)

type Movie struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Year      int       `json:"year,omitempty"`
	Runtime   Runtime   `json:"runtime,omitempty"`
	Genres    []string  `json:"genres,omitempty"`
	Version   int       `json:"version"`
}

type MovieModel struct {
	DB *sql.DB
}

func ValidateMovie(v *validator.Validator, mv *Movie) {
	v.Check(mv.Year == 0, "year", "must be provided")
	v.Check(mv.Year <= 1888, "year", "must be greater than 1888")
	v.Check(mv.Year > int(time.Now().Year()), "year", "must not be in the future")

	v.Check(mv.Runtime == 0, "runtime", "must be provided")
	v.Check(mv.Runtime < 0, "runtime", "must be a positive integer")

	v.Check(mv.Genres == nil, "genres", "must be provided")
	v.Check(len(mv.Genres) == 0, "genres", "must contain at least 1 genre")
	v.Check(len(mv.Genres) > 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.HasDuplicate(mv.Genres), "genres", "must not contain duplicate values")
}

func (m *MovieModel) Insert(movie *Movie) error {
	query := `
	INSERT INTO movies (title, year, runtime, genres)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at, version
	`

	args := []any{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)}

	return m.DB.QueryRow(query, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
}

func (m *MovieModel) Get(id int) (*Movie, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
	SELECT id, created_at, title, year, runtime, genres, version
	FROM movies
	WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var mv Movie

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&mv.ID,
		&mv.CreatedAt,
		&mv.Title,
		&mv.Year,
		&mv.Runtime,
		pq.Array(&mv.Genres),
		&mv.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &mv, nil
}

func (m *MovieModel) Update(movie *Movie) error {
	query := `
		UPDATE movies
		SET title = $1, year = $2, runtime = $3, genres = $4, version = version + 1
		WHERE id = $5 AND version = $6
		RETURNING version
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{
		movie.Title,
		movie.Year,
		movie.Runtime,
		pq.Array(movie.Genres),
		movie.ID,
		movie.Version,
	}

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&movie.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (m *MovieModel) Delete(id int) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE from movies
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (m MovieModel) GetAll(title string, genres []string, filters Filters) ([]*Movie, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, created_at, title, year, runtime, genres, version
		FROM movies
		WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (genres @> $2 OR $2 = '{}')
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4
	`, filters.Sort, filters.Order)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{title, pq.Array(genres), filters.limit(), filters.offset()}
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	movies := []*Movie{}

	for rows.Next() {
		var mv Movie

		err := rows.Scan(
			&totalRecords,
			&mv.ID,
			&mv.CreatedAt,
			&mv.Title,
			&mv.Year,
			&mv.Runtime,
			pq.Array(&mv.Genres),
			&mv.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		movies = append(movies, &mv)
	}

	err = rows.Err()
	if err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return movies, metadata, nil
}
