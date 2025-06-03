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
		r.Post("/login", h.LoginUser)
		r.Get("/verify", h.VerifyUser)
		r.Group(func(r chi.Router) {
			r.Use(h.GetCurrentUser)
			r.Put("/update-role", h.UpdateUserRole)
			r.Get("/send-verify", h.SendVerifyUser)
			r.Post("/add-editor", h.SendVerifyUser)
		})
	})

	return r
}
