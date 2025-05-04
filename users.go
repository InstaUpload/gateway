package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	pb "github.com/InstaUpload/common/api"
	common "github.com/InstaUpload/common/types"
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
	grpcResp, err := h.userClient.CreateUser(ctx, &user)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		log.Println("error creating user: ", err)
		return
	}
	log.Printf("Response: %v", grpcResp)
	resp := struct {
		message string `json:"message"`
	}{
		message: "User created successfully",
	}
	SendJsonResponse(w, http.StatusCreated, resp)
}

func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	user := pb.LoginUserRequest{
		Email:    "gpt.sahaj28@gmail.com",
		Password: "password123",
	}
	grpcResp, err := h.userClient.LoginUser(ctx, &user)
	if err != nil {
		http.Error(w, "Failed to login user", http.StatusInternalServerError)
		log.Println("error logging in user: ", err)
		return
	}
	r.Header.Set("Authorization", grpcResp.Token)
	log.Printf("Response: %v", grpcResp)
	resp := struct {
		message string `json:"message"`
	}{
		message: "User logged in successfully",
	}
	SendJsonResponse(w, http.StatusOK, resp)
}

func (h *Handler) VerifyUser(w http.ResponseWriter, r *http.Request) {
	// TODO: Get token from query string.
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Token is needed.", http.StatusBadRequest)
		log.Println("error logging in user: ")
		return
	}
	req := pb.VerifyUserRequest{
		Token: token,
	}
	grpcResp, err := h.userClient.VerifyUser(r.Context(), &req)
	if err != nil {
		if errors.Is(err, common.ErrIncorrectDataReceived) {
			http.Error(w, "Token is expired", http.StatusUnauthorized)
			return
		}
		if errors.Is(err, common.ErrDataNotFound) {
			http.Error(w, "User not found or invalid token", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to send verification to user", http.StatusInternalServerError)
		log.Println("error sending verification token to user: ", err)
		return
	}
	log.Printf("Response: %v", grpcResp)
	resp := struct {
		message string `json:"message"`
	}{
		message: "User verified",
	}
	SendJsonResponse(w, http.StatusOK, resp)
}

func (h *Handler) SendVerifyUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	req := pb.SendVerificationUserRequest{}
	grpcResp, err := h.userClient.SendVerificationUser(ctx, &req)
	if err != nil {
		http.Error(w, "Failed to send verification to user", http.StatusInternalServerError)
		log.Println("error sending verification token to user: ", err)
		return
	}
	log.Printf("Response: %v", grpcResp)
	resp := struct {
		message string `json:"message"`
	}{
		message: "Verification send.",
	}
	SendJsonResponse(w, http.StatusOK, resp)
}
