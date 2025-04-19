package service

import (
	"awesomeProject/db-service/internal/domain"
	"context"
	"github.com/google/uuid"
)

type IRatingService interface {
	CreateRating(ctx context.Context, bookID, userID uuid.UUID, request *domain.RatingRequest) (*domain.Rating, error)
	UpdateRating(ctx context.Context, bookID, userID uuid.UUID, request *domain.RatingRequest) (*domain.Rating, error)
	GetRating(ctx context.Context, bookID, userID uuid.UUID) (*domain.Rating, error)
	GetBookRatings(ctx context.Context, bookID uuid.UUID) ([]*domain.Rating, error)
	GetBookStats(ctx context.Context, bookID uuid.UUID) (*domain.BookRatingStats, error)
	GetAllBookStats(ctx context.Context) ([]*domain.BookRatingStats, error)
	DeleteRating(ctx context.Context, bookID, userID uuid.UUID) error
	HandleNewBook(ctx context.Context, bookID uuid.UUID) error
}
