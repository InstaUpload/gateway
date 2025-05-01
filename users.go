package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	pb "github.com/InstaUpload/common/api"
)

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	decoder := json.NewDecoder(r.Body)
	var user pb.CreateUserRequest
	if err := decoder.Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Println("error decoding request: ", err)
		return
	}
	resp, err := h.userClient.CreateUser(ctx, &user)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		log.Println("error creating user: ", err)
		return
	}
	log.Printf("Response: %v", resp)
	resp = struct {
		message string `json:"message"`
	}{
		message: "User created successfully",
	}
	SendJsonResponse(&w, http.StatusCreated, resp)
}

func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	user := pb.LoginUserRequest{
		Email:    "gpt.sahaj28@gmail.com",
		Password: "password123",
	}
	resp, err := h.userClient.LoginUser(ctx, &user)
	if err != nil {
		http.Error(w, "Failed to login user", http.StatusInternalServerError)
		log.Println("error logging in user: ", err)
		return
	}
	r.Header().Set("Authorization", resp.Token)
	log.Printf("Response: %v", resp)
	resp = struct {
		message string `json:"message"`
	}{
		message: "User logged in successfully",
	}
	SendJsonResponse(&w, http.StatusOK, resp)
}

func (h *Handler) SendVerifyUser(w http.ResponseWriter, r *http.Request) {
	errors.New("not implemented")
}
