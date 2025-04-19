package repository

import (
	"awesomeProject/db-service/internal/domain"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	ErrRatingNotFound      = errors.New("rating not found")
	ErrRatingAlreadyExists = errors.New("rating already exists for this book and user")
)

type RatingRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewRatingRepository(db *pgxpool.Pool) IRatingRepository {
	return &RatingRepositoryImpl{
		db: db,
	}
}

func (r *RatingRepositoryImpl) Create(ctx context.Context, rating *domain.Rating) error {
	if rating.ID == uuid.Nil {
		rating.ID = uuid.New()
	}

	query := `
		INSERT INTO ratings (id, book_id, user_id, rating, comment, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Exec(ctx, query,
		rating.ID, rating.BookID, rating.UserID, rating.Rating, rating.Comment, rating.CreatedAt)
	if err != nil {
		return fmt.Errorf("error creating rating in database: %w", err)
	}
	return nil
}

func (r *RatingRepositoryImpl) Update(ctx context.Context, rating *domain.Rating) error {
	query := `
		UPDATE ratings
		SET rating = $1, comment = $2
		WHERE book_id = $3 AND user_id = $4
		RETURNING id
	`
	var id uuid.UUID
	err := r.db.QueryRow(ctx, query,
		rating.Rating, rating.Comment, rating.BookID, rating.UserID).Scan(&id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrRatingNotFound
		}
		return fmt.Errorf("error updating rating in database: %w", err)
	}

	return nil
}

func (r *RatingRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*domain.Rating, error) {
	var rating domain.Rating
	query := `
		SELECT id, book_id, user_id, rating, comment, created_at
		FROM ratings
		WHERE id = $1
	`

	err := r.db.QueryRow(ctx, query, id).Scan(
		&rating.ID, &rating.BookID, &rating.UserID, &rating.Rating, &rating.Comment, &rating.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("rating with ID %s not found", id)
		}
		return nil, fmt.Errorf("error requesting rating with ID %s: %w", id, err)
	}

	return &rating, nil
}

func (r *RatingRepositoryImpl) GetByBookAndUser(ctx context.Context, bookID, userID uuid.UUID) (*domain.Rating, error) {
	var rating domain.Rating
	query := `
		SELECT id, book_id, user_id, rating, comment, created_at
		FROM ratings
		WHERE book_id = $1 AND user_id = $2
	`

	err := r.db.QueryRow(ctx, query, bookID, userID).Scan(
		&rating.ID, &rating.BookID, &rating.UserID, &rating.Rating, &rating.Comment, &rating.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("rating for book %s and user %s not found", bookID, userID)
		}
		return nil, fmt.Errorf("error getting rating from database: %w", err)
	}

	return &rating, nil
}

func (r *RatingRepositoryImpl) GetByBookID(ctx context.Context, bookID uuid.UUID) ([]*domain.Rating, error) {
	query := `
		SELECT id, book_id, user_id, rating, comment, created_at
		FROM ratings
		WHERE book_id = $1
		ORDER BY rating DESC
	`

	rows, err := r.db.Query(ctx, query, bookID)
	if err != nil {
		return nil, fmt.Errorf("error getting ratings for book from database: %w", err)
	}
	defer rows.Close()

	ratings := []*domain.Rating{}

	for rows.Next() {
		var rating domain.Rating
		err := rows.Scan(&rating.ID, &rating.BookID, &rating.UserID, &rating.Rating, &rating.Comment, &rating.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning rating results: %w", err)
		}
		ratings = append(ratings, &rating)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after processing rating results: %w", err)
	}

	return ratings, nil
}

func (r *RatingRepositoryImpl) GetBookStats(ctx context.Context, bookID uuid.UUID) (*domain.BookRatingStats, error) {
	var stats domain.BookRatingStats
	query := `
		SELECT book_id, average_rating, rating_count
		FROM book_average_ratings
		WHERE book_id = $1
	`

	err := r.db.QueryRow(ctx, query, bookID).Scan(
		&stats.BookID, &stats.AverageRating, &stats.RatingCount)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("stats for book %s not found", bookID)
		}
		return nil, fmt.Errorf("error getting book rating stats from database: %w", err)
	}

	return &stats, nil
}

func (r *RatingRepositoryImpl) GetAllBookStats(ctx context.Context) ([]*domain.BookRatingStats, error) {
	query := `
		SELECT book_id, average_rating, rating_count
		FROM book_average_ratings
		ORDER BY average_rating DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error getting all book rating stats from database: %w", err)
	}
	defer rows.Close()

	stats := []*domain.BookRatingStats{}

	for rows.Next() {
		var stat domain.BookRatingStats
		err := rows.Scan(&stat.BookID, &stat.AverageRating, &stat.RatingCount)
		if err != nil {
			return nil, fmt.Errorf("error scanning book stats results: %w", err)
		}
		stats = append(stats, &stat)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after processing book stats results: %w", err)
	}

	return stats, nil
}

func (r *RatingRepositoryImpl) Delete(ctx context.Context, bookID, userID uuid.UUID) error {
	query := `
		DELETE FROM ratings
		WHERE book_id = $1 AND user_id = $2
	`

	commandTag, err := r.db.Exec(ctx, query, bookID, userID)
	if err != nil {
		return fmt.Errorf("error deleting rating from database: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return ErrRatingNotFound
	}

	return nil
}

func (r *RatingRepositoryImpl) InitEmptyRating(ctx context.Context, bookID uuid.UUID) error {
	return nil
}
