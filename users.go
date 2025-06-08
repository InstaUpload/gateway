package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	pb "github.com/InstaUpload/common/api"
	common "github.com/InstaUpload/common/types"
	"github.com/go-chi/chi"
)

type MessageResponse struct {
	Message string `json:"message"`
}

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// CreateUser godoc
//
//	@Summary		Create User
//	@Description	Create a new user
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			user	body		CreateUserRequest	true	"User details"
//	@Success		201		{object}	MessageResponse
//	@Failure		400		{object}	MessageResponse
//	@Failure		500		{object}	MessageResponse
//	@Router			/v1/users/create [post]
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
	resp := MessageResponse{
		Message: "User created successfully",
	}
	SendJsonResponse(w, http.StatusCreated, resp)
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginUser godoc
//
//	@Summary		Login User
//	@Description	Login to an existing user
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			user	body		LoginUserRequest	true	"User login details"
//	@Success		200		{object}	MessageResponse
//	@Failure		500		{object}	MessageResponse
//	@Router			/v1/users/login [post]
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
	r.Header.Set("Authorization", "Bearer "+grpcResp.Token)
	log.Printf("Response: %v", grpcResp)
	SendJsonResponse(w, http.StatusOK, grpcResp)
}

// VerifyUser godoc
//
//	@Summary		Verify User
//	@Description	Verify a existing user
//	@Tags			Users
//	@Produce		json
//	@Param			token	query		string	true	"Token send to user's mail for verification"
//	@Success		200		{object}	MessageResponse
//	@Failure		400		{object}	MessageResponse
//	@Failure		401		{object}	MessageResponse
//	@Failure		404		{object}	MessageResponse
//	@Failure		500		{object}	MessageResponse
//	@Router			/v1/users/verify [get]
func (h *Handler) VerifyUser(w http.ResponseWriter, r *http.Request) {
	// Get token from query string.
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
		Message string `json:"message"`
	}{
		Message: "User verified",
	}
	SendJsonResponse(w, http.StatusOK, resp)
}

// SendVerifyUser godoc
//
//	@Summary		Send Verify User
//	@Description	Send Verify token to a existing user
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			token	query		string	true	"Token send to user's mail for verification"
//	@Success		200		{object}	MessageResponse
//	@Failure		500		{object}	MessageResponse
//	@Security		ApiKeyAuth
//	@Router			/v1/users/send-verify [get]
func (h *Handler) SendVerifyUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Get Current user from ctx, and pass it in SendVerificationUserRequest.
	req := pb.SendVerificationUserRequest{
		CurrentUser: ctx.Value(common.CurrentUserKey).(*pb.AuthUserResponse),
	}
	grpcResp, err := h.userClient.SendVerificationUser(ctx, &req)
	if err != nil {
		http.Error(w, "Failed to send verification to user", http.StatusInternalServerError)
		log.Println("error sending verification token to user: ", err)
		return
	}
	log.Printf("Response: %v", grpcResp)
	resp := struct {
		Message string `json:"message"`
	}{
		Message: "Verification send.",
	}
	SendJsonResponse(w, http.StatusOK, resp)
}

type UpdateUserRoleRequest struct {
	UserID   string `json:"user_id"`
	RoleName string `json:"role_name"`
}

// UpdateUserRole godoc
//
//	@Summary		Update User Role
//	@Description	Update the role of an existing user
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			data	body		UpdateUserRoleRequest	true	"User ID and role name"
//	@Success		200		{object}	MessageResponse
//	@Failure		400		{object}	MessageResponse
//	@Failure		500		{object}	MessageResponse
//	@Security		ApiKeyAuth
//	@Router			/v1/users/update-role [put]
func (h *Handler) UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	// get user id and role name from request body.
	decoder := json.NewDecoder(r.Body)
	var req pb.UpdateUserRoleRequest
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Println("error decoding request: ", err)
		return
	}
	// Get Current user from ctx, and pass it in UpdateUserRoleRequest.
	req.CurrentUser = ctx.Value(common.CurrentUserKey).(*pb.AuthUserResponse)
	grpcResp, err := h.userClient.UpdateUserRole(ctx, &req)
	if err != nil {
		http.Error(w, "Failed to update user role", http.StatusInternalServerError)
		log.Println("error updating user role: ", err)
		return
	}
	log.Printf("Response: %v", grpcResp)
	resp := MessageResponse{
		Message: "User role updated successfully",
	}
	SendJsonResponse(w, http.StatusOK, resp)
}

// AddEditorUser godoc
//
//	@Summary		Add Editor User
//	@Description	Add a user as an editor
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			token	query		string	true	"Token for adding editor user"
//	@Success		200		{object}	MessageResponse
//	@Failure		400		{object}	MessageResponse
//	@Failure		500		{object}	MessageResponse
//	@Security		ApiKeyAuth
//	@Router			/v1/users/add-editor [post]
func (h *Handler) AddEditorUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Token is needed.", http.StatusBadRequest)
		log.Println("Token not provided")
		return
	}
	req := pb.AddEditorUserRequest{
		Token: token,
	}
	_, err := h.userClient.AddEditorUser(ctx, &req)
	if err != nil {
		http.Error(w, "Failed to add editor user", http.StatusInternalServerError)
		log.Println("error adding editor user: ", err)
		return
	}
	resp := MessageResponse{
		Message: "Added editor user successfully",
	}
	SendJsonResponse(w, http.StatusOK, resp)
}

// SendEditorInvite godoc
//
//	@Summary		Send Editor Invite
//	@Description	Send an invite to a user to become an editor
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			u	path		int64	true	"User ID to send editor invite"
//	@Success		200	{object}	MessageResponse
//	@Failure		400	{object}	MessageResponse
//	@Failure		500	{object}	MessageResponse
//	@Security		ApiKeyAuth
//	@Router			/v1/users/send-editor-invite/{u} [put]
func (h *Handler) SendEditorInvite(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	uId := chi.URLParam(r, "u")
	userId, err := strconv.ParseInt(uId, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		log.Println("error parsing user ID: ", err)
		return
	}
	req := pb.SendEditorUserRequest{
		UserId:      userId,
		CurrentUser: ctx.Value(common.CurrentUserKey).(*pb.AuthUserResponse),
	}
	_, err = h.userClient.SendEditorUser(ctx, &req)
	if err != nil {
		http.Error(w, "Failed to send editor invite", http.StatusInternalServerError)
		log.Println("error sending editor invite: ", err)
		return
	}
	resp := MessageResponse{
		Message: "Sent editor invite successfully",
	}
	SendJsonResponse(w, http.StatusOK, resp)
}

type ResetUserPasswordRequest struct {
	Email string `json:"email"`
}

// ResetUserPassword godoc
//
//	@Summary		Reset User Password
//	@Description	Reset the password of an existing user
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			data	body		ResetUserPasswordRequest	true	"User email to reset password"
//	@Success		200		{object}	MessageResponse
//	@Failure		400		{object}	MessageResponse
//	@Failure		500		{object}	MessageResponse
//	@Router			/v1/users/reset-password [post]
func (h *Handler) ResetUserPassword(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var req pb.ResetUserPasswordRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Println("error decoding request: ", err)
		return
	}

	grpcResp, err := h.userClient.ResetUserPassword(ctx, &req)
	if err != nil {
		http.Error(w, "Failed to reset user password", http.StatusInternalServerError)
		log.Println("error resetting user password: ", err)
		return
	}
	log.Printf("Response: %v", grpcResp)

	resp := MessageResponse{
		Message: "Password reset successfully",
	}
	SendJsonResponse(w, http.StatusOK, resp)
}

type UpdateUserPasswordRequest struct {
	Password string `json:"password"`
}

// UpdateUserPassword godoc
//
//	@Summary		Update User Password
//	@Description	Update the password of an existing user
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			token	query		string						true	"Token for updating user password"
//	@Param			data	body		UpdateUserPasswordRequest	true	"New password for the user"
//	@Success		200		{object}	MessageResponse
//	@Failure		400		{object}	MessageResponse
//	@Failure		500		{object}	MessageResponse
//	@Router			/v1/users/update-password [post]
func (h *Handler) UpdateUserPassword(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	// Get token from query string.
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Token is needed.", http.StatusBadRequest)
		log.Println("Token not provided")
		return
	}

	var req pb.UpdateUserPasswordRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Println("error decoding request: ", err)
		return
	}
	req.Token = token

	grpcResp, err := h.userClient.UpdateUserPassword(ctx, &req)
	if err != nil {
		http.Error(w, "Failed to update user password", http.StatusInternalServerError)
		log.Println("error updating user password: ", err)
		return
	}
	log.Printf("Response: %v", grpcResp)

	resp := MessageResponse{
		Message: "Password updated successfully",
	}
	SendJsonResponse(w, http.StatusOK, resp)
}
