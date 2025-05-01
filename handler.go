package main

import (
	"net/http"
	"time"

	pb "github.com/InstaUpload/common/api"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Handler struct {
	userClient pb.UserServiceClient
}

func (h *Handler) mount() http.Handler {
	r := chi.NewRouter()
	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Route("/users", func(r chi.Router) {
		r.Post("/create", h.CreateUser)
	})

	return r
}
