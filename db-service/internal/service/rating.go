package service

import (
	"awesomeProject/db-service/internal/domain"
	"awesomeProject/db-service/internal/repository"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type RatingServiceImpl struct {
	ratingRepo repository.IRatingRepository
}

func NewRatingService(ratingRepo repository.IRatingRepository) IRatingService {
	return &RatingServiceImpl{
		ratingRepo: ratingRepo,
	}
}

func (s *RatingServiceImpl) CreateRating(ctx context.Context, bookID, userID uuid.UUID, request *domain.RatingRequest) (*domain.Rating, error) {

	_, err := s.ratingRepo.GetByBookAndUser(ctx, bookID, userID)
	if err == nil {

		return nil, repository.ErrRatingAlreadyExists
	}

	rating := &domain.Rating{
		ID:        uuid.New(),
		BookID:    bookID,
		UserID:    userID,
		Rating:    request.Rating,
		Comment:   request.Comment,
		CreatedAt: time.Now(),
	}

	if err := s.ratingRepo.Create(ctx, rating); err != nil {
		return nil, fmt.Errorf("failed to create rating: %w", err)
	}

	return rating, nil
}

func (s *RatingServiceImpl) UpdateRating(ctx context.Context, bookID, userID uuid.UUID, request *domain.RatingRequest) (*domain.Rating, error) {

	existingRating, err := s.ratingRepo.GetByBookAndUser(ctx, bookID, userID)
	if err != nil {
		return nil, fmt.Errorf("rating not found: %w", err)
	}

	existingRating.Rating = request.Rating
	existingRating.Comment = request.Comment

	if err := s.ratingRepo.Update(ctx, existingRating); err != nil {
		return nil, fmt.Errorf("failed to update rating: %w", err)
	}

	return existingRating, nil
}

func (s *RatingServiceImpl) GetRating(ctx context.Context, bookID, userID uuid.UUID) (*domain.Rating, error) {
	rating, err := s.ratingRepo.GetByBookAndUser(ctx, bookID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get rating: %w", err)
	}
	return rating, nil
}

func (s *RatingServiceImpl) GetBookRatings(ctx context.Context, bookID uuid.UUID) ([]*domain.Rating, error) {
	ratings, err := s.ratingRepo.GetByBookID(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("failed to get book ratings: %w", err)
	}
	return ratings, nil
}

func (s *RatingServiceImpl) GetBookStats(ctx context.Context, bookID uuid.UUID) (*domain.BookRatingStats, error) {
	stats, err := s.ratingRepo.GetBookStats(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("failed to get book stats: %w", err)
	}
	return stats, nil
}

func (s *RatingServiceImpl) GetAllBookStats(ctx context.Context) ([]*domain.BookRatingStats, error) {
	stats, err := s.ratingRepo.GetAllBookStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all book stats: %w", err)
	}
	return stats, nil
}

func (s *RatingServiceImpl) DeleteRating(ctx context.Context, bookID, userID uuid.UUID) error {
	if err := s.ratingRepo.Delete(ctx, bookID, userID); err != nil {
		return fmt.Errorf("failed to delete rating: %w", err)
	}
	return nil
}

func (s *RatingServiceImpl) HandleNewBook(ctx context.Context, bookID uuid.UUID) error {

	if err := s.ratingRepo.InitEmptyRating(ctx, bookID); err != nil {
		return fmt.Errorf("failed to initialize empty rating for new book: %w", err)
	}
	return nil
}
