package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Movie struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type MovieRepo struct {
	pool *pgxpool.Pool
}

func NewMovieRepo(pool *pgxpool.Pool) *MovieRepo {
	return &MovieRepo{
		pool: pool,
	}
}

func (repo *MovieRepo) GetAllMovies() ([]Movie, error) {

	conn, err := repo.pool.Acquire(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to acquire database connection")
	}
	defer conn.Release()

	query := "SELECT id, title FROM movies"
	rows, err := conn.Query(context.Background(), query)

	var movies []Movie
	for rows.Next() {
		var movie Movie
		err := rows.Scan(&movie.ID, &movie.Title)
		if err != nil {
			return nil, fmt.Errorf("error scanning movie row %v", err)
		}
		movies = append(movies, movie)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error in rows: %v", err)
	}

	return movies, nil
}
