package domain

import (
	"time"

	"github.com/google/uuid"
)

type Rating struct {
	ID        uuid.UUID `json:"id" db:"id"`
	BookID    uuid.UUID `json:"book_id" db:"book_id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Rating    int       `json:"rating" db:"rating"`
	Comment   string    `json:"comment,omitempty" db:"comment"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type BookRatingStats struct {
	BookID        uuid.UUID `json:"book_id" db:"book_id"`
	AverageRating float64   `json:"average_rating" db:"average_rating"`
	RatingCount   int       `json:"rating_count" db:"rating_count"`
}

type RatingRequest struct {
	Rating  int    `json:"rating" binding:"required,min=0,max=10"`
	Comment string `json:"comment"`
}
