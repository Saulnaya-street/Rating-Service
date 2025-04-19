package handler

import (
	"awesomeProject/db-service/internal/domain"
	"awesomeProject/db-service/internal/repository"
	"awesomeProject/db-service/internal/service"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type IRatingHandler interface {
	GetBookRatings(w http.ResponseWriter, r *http.Request)
	GetBookStats(w http.ResponseWriter, r *http.Request)
	GetAllBookStats(w http.ResponseWriter, r *http.Request)
	CreateRating(w http.ResponseWriter, r *http.Request)
	GetRating(w http.ResponseWriter, r *http.Request)
	UpdateRating(w http.ResponseWriter, r *http.Request)
	DeleteRating(w http.ResponseWriter, r *http.Request)
}

type RatingHandler struct {
	ratingService service.IRatingService
}

func NewRatingHandler(ratingService service.IRatingService) IRatingHandler {
	return &RatingHandler{
		ratingService: ratingService,
	}
}

func (h *RatingHandler) GetBookRatings(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID, err := uuid.Parse(vars["bookID"])
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	ratings, err := h.ratingService.GetBookRatings(r.Context(), bookID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ratings)
}

func (h *RatingHandler) GetBookStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID, err := uuid.Parse(vars["bookID"])
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	stats, err := h.ratingService.GetBookStats(r.Context(), bookID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func (h *RatingHandler) GetAllBookStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.ratingService.GetAllBookStats(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func (h *RatingHandler) CreateRating(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID, err := uuid.Parse(vars["bookID"])
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(vars["userID"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var request domain.RatingRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rating, err := h.ratingService.CreateRating(r.Context(), bookID, userID, &request)
	if err != nil {
		if err == repository.ErrRatingAlreadyExists {
			http.Error(w, "Rating already exists for this book and user", http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(rating)
}

func (h *RatingHandler) GetRating(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID, err := uuid.Parse(vars["bookID"])
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(vars["userID"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	rating, err := h.ratingService.GetRating(r.Context(), bookID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rating)
}

func (h *RatingHandler) UpdateRating(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID, err := uuid.Parse(vars["bookID"])
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(vars["userID"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var request domain.RatingRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rating, err := h.ratingService.UpdateRating(r.Context(), bookID, userID, &request)
	if err != nil {
		if err == repository.ErrRatingNotFound {
			http.Error(w, "Rating not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rating)
}

func (h *RatingHandler) DeleteRating(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID, err := uuid.Parse(vars["bookID"])
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(vars["userID"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if err := h.ratingService.DeleteRating(r.Context(), bookID, userID); err != nil {
		if err == repository.ErrRatingNotFound {
			http.Error(w, "Rating not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
