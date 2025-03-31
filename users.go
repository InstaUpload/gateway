package main

import (
	"context"
	"log"
	"net/http"
	"time"

	pb "github.com/InstaUpload/common/api"
)

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	user := pb.CreateUserRequest{
		Name:     "Sahaj",
		Email:    "gpt.sahaj28@gmail.com",
		Password: "password123",
	}
	resp, err := h.userClient.CreateUser(ctx, &user)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		log.Println("error creating user: ", err)
		return
	}

	log.Printf("Response: %v", resp)
	w.Write([]byte("User created successfully"))
}
