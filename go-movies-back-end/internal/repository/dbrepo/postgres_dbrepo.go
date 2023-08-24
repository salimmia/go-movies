package dbrepo

import (
	"backend/internal/models"
	"context"
	"database/sql"
	"log"
	"time"
)

type PostgresDBRepo struct {
	DB *sql.DB
}

const dbTimeout = time.Second * 3

func (m *PostgresDBRepo) Connection() *sql.DB {
	return m.DB
}

func (m *PostgresDBRepo) AllMovies() ([]*models.Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		select
			id, title, release_date, runtime,
			mpaa_rating, description, coalesce(image, ''),
			created_at, updated_at
		from
			movies
		order by
			title
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []*models.Movie

	for rows.Next() {
		var movie models.Movie
		err := rows.Scan(
			&movie.ID,
			&movie.Title,
			&movie.ReleaseDate,
			&movie.RunTime,
			&movie.MPAARating,
			&movie.Description,
			&movie.Image,
			&movie.CreatedAt,
			&movie.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		movies = append(movies, &movie)
	}

	return movies, nil
}

func (m *PostgresDBRepo) GetUserByEmail(email string) (*models.User, error){
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password,
		created_at, updated_at from users where email = $1`
		
	var user models.User

	row := m.DB.QueryRowContext(ctx, query, email)

	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil{
		return nil, err
	}
	return &user, nil
}

func (m *PostgresDBRepo) GetUserByID(id int) (*models.User, error){
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password,
		created_at, updated_at from users where id = $1`
		
	var user models.User

	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil{
		return nil, err
	}
	return &user, nil
}

func (m *PostgresDBRepo) OneMovie(id int) (*models.Movie, error){
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, title, release_date, runtime, mpaa_rating, description, coalesce(image, ''), created_at, updated_at from movies where id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	var movie models.Movie

	err := row.Scan(
		 &movie.ID,
		 &movie.Title,
		 &movie.RunTime,
		 &movie.MPAARating,
		 &movie.Description,
		 &movie.Image,
		 &movie.CreatedAt,
		 &movie.UpdatedAt,
	)

	if err != nil{
		return &movie, err
	}

	query = `select g.id, g.genre from movies_genres mg
		left join genres g on(mg.genre_id = g.id)
		where mg.movie_id = $1
		order by g.genre`

	rows, err := m.DB.QueryContext(ctx, query, id)

	if err != nil && err != sql.ErrNoRows{
		return nil, err
	}
	defer rows.Close()

	var genres []*models.Genre

	for rows.Next(){
		var genre models.Genre;

		err := rows.Scan(
			&genre.ID,
			&genre.Genre,
		)

		if err != nil{
			return nil, err
		}
		genres = append(genres, &genre)
	}

	movie.Genres = genres

	log.Println(movie)

	return &movie, nil
}

func (m *PostgresDBRepo) OneMovieForEdit(id int) (*models.Movie, []*models.Genre, error){
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, title, release_date, runtime, mpaa_rating, description, coalesce(image, ''), created_at, updated_at from movies where id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	var movie models.Movie

	err := row.Scan(
		 &movie.ID,
		 &movie.Title,
		 &movie.RunTime,
		 &movie.MPAARating,
		 &movie.Description,
		 &movie.Image,
		 &movie.CreatedAt,
		 &movie.UpdatedAt,
	)

	if err != nil{
		return &movie, nil, err
	}

	query = `select g.id, g.genre from movies_genres mg
		left join genres g on(mg.genre_id = g.id)
		where mg.movie_id = $1
		order by g.genre`

	rows, err := m.DB.QueryContext(ctx, query, id)

	if err != nil && err != sql.ErrNoRows{
		return nil, nil, err
	}
	defer rows.Close()

	var genres []*models.Genre
	var genresArray []int

	for rows.Next(){
		var genre models.Genre;

		err := rows.Scan(
			&genre.ID,
			&genre.Genre,
		)

		if err != nil{
			return nil, nil, err
		}
		genres = append(genres, &genre)
		genresArray = append(genresArray, genre.ID)
	}

	movie.Genres = genres
	movie.GenresArray = genresArray

	log.Println(movie)

	var allGenre []*models.Genre
	
	query = `
		select id, genre from genres order by genre	
	`

	gRows, err := m.DB.QueryContext(ctx, query)

	if err != nil{
		return nil, nil, err;
	}
	gRows.Close()

	for gRows.Next(){
		var g models.Genre

		err = gRows.Scan(
			&g.ID,
			&g.Genre,
		)
		if err != nil{
			return nil, nil, err
		}
		allGenre = append(allGenre, &g)
	}

	return &movie, allGenre, nil
}