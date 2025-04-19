package repository

import (
	"awesomeProject/db-service/internal/domain"
	"context"
	"github.com/google/uuid"
)

type IRatingRepository interface {
	Create(ctx context.Context, rating *domain.Rating) error
	Update(ctx context.Context, rating *domain.Rating) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Rating, error)
	GetByBookAndUser(ctx context.Context, bookID, userID uuid.UUID) (*domain.Rating, error)
	GetByBookID(ctx context.Context, bookID uuid.UUID) ([]*domain.Rating, error)
	GetBookStats(ctx context.Context, bookID uuid.UUID) (*domain.BookRatingStats, error)
	GetAllBookStats(ctx context.Context) ([]*domain.BookRatingStats, error)
	Delete(ctx context.Context, bookID, userID uuid.UUID) error
	InitEmptyRating(ctx context.Context, bookID uuid.UUID) error
}
