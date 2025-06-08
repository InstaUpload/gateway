package main

import (
	"net/http"
	"time"

	pb "github.com/InstaUpload/common/api"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
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
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:5000/swagger/doc.json"), //The url pointing to API definition
	))
	r.Route("/v1", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Post("/create", h.CreateUser)
			r.Post("/login", h.LoginUser)
			r.Get("/verify", h.VerifyUser)
			r.Post("/reset-password", h.ResetUserPassword)
			r.Post("/update-password", h.UpdateUserPassword)
			r.Group(func(r chi.Router) {
				r.Use(h.GetCurrentUser)
				r.Put("/update-role", h.UpdateUserRole)
				r.Get("/send-verify", h.SendVerifyUser)
				r.Put("/add-editor", h.AddEditorUser)
				r.Put("/send-editor-invite/{u}", h.SendEditorInvite)
			})
		})
	})

	return r
}
