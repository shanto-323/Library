package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	userservice "github.com/shanto-323/Library/user_service"
)

func (s *Server) SignUpHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if r.Method != http.MethodPost {
		return fmt.Errorf("invalid mathod")
	}

	if r.Body == nil {
		http.Error(w, "empty body", http.StatusBadRequest)
		return nil
	}
	defer r.Body.Close()

	user := &userservice.UserModel{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return err
	}

	if user.UserType != "ADMIN" && user.UserType != "USER" {
		WriteJson(w, http.StatusBadRequest, "invalid user type")
		return nil
	}

	if user.UserType == "ADMIN" && user.Phone == "" {
		WriteJson(w, http.StatusBadRequest, "phone number requird")
		return nil
	}
	resp, err := s.userServiceClient.SignUp(ctx, user.Name, user.Password, user.Email, user.Phone, user.UserType)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "access_token",
		Value: resp.Token,
		Path:  "/",
	})
	return WriteJson(w, http.StatusOK, resp)
}

func (s *Server) SignInHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if r.Method != http.MethodPost {
		return fmt.Errorf("invalid mathod")
	}

	if r.Body == nil {
		http.Error(w, "empty body", http.StatusBadRequest)
		return nil
	}
	defer r.Body.Close()

	user := &userservice.UserModel{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return err
	}

	if user.Email == "" || user.Password == "" {
		return fmt.Errorf("empty field")
	}

	resp, err := s.userServiceClient.SignIn(ctx, user.Email, user.Password)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "access_token",
		Value: resp.Token,
		Path:  "/",
	})
	return WriteJson(w, http.StatusOK, resp)
}

func (s *Server) SignOutHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if r.Method != http.MethodPost {
		return fmt.Errorf("invalid mathod")
	}
	id := mux.Vars(r)["id"]

	resp, err := s.userServiceClient.Logout(ctx, id)
	if err != nil {
		return err
	}
	return WriteJson(w, http.StatusOK, resp)
}

func (s *Server) UpdateUserHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if r.Method != http.MethodPatch {
		return fmt.Errorf("invalid mathod")
	}
	// small issue ,anyone can send request with id
	// number even though id not match with user.
	id := mux.Vars(r)["id"]

	if r.Body == nil {
		http.Error(w, "empty body", http.StatusBadRequest)
		return nil
	}
	defer r.Body.Close()

	user := &userservice.UserModel{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return err
	}
	user.ID = id

	resp, err := s.userServiceClient.UpdateUser(ctx, user)
	if err != nil {
		return err
	}

	return WriteJson(w, http.StatusOK, *resp)
}

func (s *Server) DeleteUserHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if r.Method != http.MethodDelete {
		return fmt.Errorf("invalid mathod")
	}

	id := mux.Vars(r)["id"]
	resp, err := s.userServiceClient.DeleteUser(ctx, id)
	if err != nil {
		return err
	}

	return WriteJson(w, http.StatusOK, *resp)
}

func (s *Server) GetALlUserHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if r.Method != http.MethodGet {
		return fmt.Errorf("invalid mathod")
	}

	u_type := r.Context().Value(userTypeKey).(string)
	if u_type != "ADMIN" {
		return fmt.Errorf("not allowed")
	}

	limit, err := perseInt("limit", r)
	if err != nil {
		return err
	}
	offset, err := perseInt("offset", r)
	if err != nil {
		return err
	}
	users, err := s.userServiceClient.GetUsers(ctx, limit, offset)
	if err != nil {
		return err
	}

	return WriteJson(w, http.StatusOK, users)
}

func (s *Server) NewRefreshTokenHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if r.Method != http.MethodGet {
		return fmt.Errorf("invalid mathod")
	}

	r_token := r.Header.Get("x-jwt-token")
	if r_token == "" {
		return fmt.Errorf("nil token")
	}

	_, err := ValidateToken(r_token)
	if err != nil {
		return err
	}
	id := mux.Vars(r)["id"]
	newToken, err := s.userServiceClient.GetAccessToken(ctx, id, r_token)
	if err != nil {
		return err
	}
	http.SetCookie(w, &http.Cookie{
		Name:  "access_token",
		Value: *newToken,
		Path:  "/",
	})
	return WriteJson(w, http.StatusOK, newToken)
}
