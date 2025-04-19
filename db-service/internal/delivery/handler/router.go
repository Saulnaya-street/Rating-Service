package handler

import (
	"github.com/gorilla/mux"
)

type Router struct {
	ratingHandler IRatingHandler
}

func NewRouter(ratingHandler IRatingHandler) *Router {
	return &Router{
		ratingHandler: ratingHandler,
	}
}

func (r *Router) RegisterRoutes(router *mux.Router) {

	router.HandleFunc("/api/ratings/stats", r.ratingHandler.GetAllBookStats).Methods("GET")

	router.HandleFunc("/api/books/{bookID}/ratings", r.ratingHandler.GetBookRatings).Methods("GET")
	router.HandleFunc("/api/books/{bookID}/ratings/stats", r.ratingHandler.GetBookStats).Methods("GET")

	router.HandleFunc("/api/books/{bookID}/users/{userID}/ratings", r.ratingHandler.CreateRating).Methods("POST")
	router.HandleFunc("/api/books/{bookID}/users/{userID}/ratings", r.ratingHandler.GetRating).Methods("GET")
	router.HandleFunc("/api/books/{bookID}/users/{userID}/ratings", r.ratingHandler.UpdateRating).Methods("PUT")
	router.HandleFunc("/api/books/{bookID}/users/{userID}/ratings", r.ratingHandler.DeleteRating).Methods("DELETE")
}
