package main

import (
	"context"
	"errors"
	"log"
	"net/http"

	pb "github.com/InstaUpload/common/api"
	common "github.com/InstaUpload/common/types"
)

func (h *Handler) GetCurrentUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the token from the request header
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		var req = pb.AuthUserRequest{}
		req.Token = token
		resp, err := h.userClient.AuthUser(r.Context(), &req)
		if err != nil {
			if errors.Is(err, common.ErrIncorrectDataReceived) {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			if errors.Is(err, common.ErrDataNotFound) {
				http.Error(w, "User not found", http.StatusNotFound)
				return
			}
			log.Println("error authenticating user: ", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		// Set resp(User) in the request context
		ctx := context.WithValue(r.Context(), common.CurrentUserKey, resp)
		// Call the next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
