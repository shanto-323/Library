package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/shanto-323/Library/books"
	userservice "github.com/shanto-323/Library/user_service"
)

type Server struct {
	IpAddr            string
	bookClient        *books.Client
	userServiceClient *userservice.Client
}

func NewServer(ipAddr string, bookClientUrl string, userServiceClientUrl string) (*Server, error) {
	bookClient, err := books.NewClient(bookClientUrl)
	if err != nil {
		log.Println("NewServer Error", err)
		return nil, err
	}

	userServiceClient, err := userservice.NewClient(userServiceClientUrl)
	if err != nil {
		log.Println("NewServer Error", err)
		return nil, err
	}

	return &Server{
		IpAddr:            ipAddr,
		bookClient:        bookClient,
		userServiceClient: userServiceClient,
	}, nil
}

func (s *Server) Start() error {
	r := mux.NewRouter()
	router := r.PathPrefix("/library/v2").Subrouter()

	//books
	bookRouter := router.PathPrefix("/books").Subrouter()
	bookRouter.Use(JwtMiddleWere)

	bookRouter.HandleFunc("/search", createHandlerFunc(s.SearchBookHandler)).Methods("GET")
	bookRouter.HandleFunc("", createHandlerFunc(s.CreateBookHandler)).Methods("POST")
	bookRouter.HandleFunc("/{isbn}", createHandlerFunc(s.UpdateBookHandler)).Methods("PATCH")
	bookRouter.HandleFunc("/{isbn}", createHandlerFunc(s.DeleteBookHandler)).Methods("DELETE")
	bookRouter.HandleFunc("/{isbn}", createHandlerFunc(s.GetBookHandler)).Methods("GET")
	bookRouter.HandleFunc("", createHandlerFunc(s.GetAllBooksHandler)).Methods("GET")

	//user
	userServiceRouter := router.PathPrefix("/user").Subrouter()
	userServiceRouter.HandleFunc("/signup", createHandlerFunc(s.SignUpHandler)).Methods("POST")
	userServiceRouter.HandleFunc("/login", createHandlerFunc(s.SignInHandler)).Methods("GET")
	userServiceRouter.HandleFunc("/logout/{id}", createHandlerFunc(s.SignOutHandler)).Methods("POST")
	userServiceRouter.HandleFunc("/token", createHandlerFunc(s.NewRefreshTokenHandler)).Methods("POST")

	fmt.Println("Api running.. ")
	return http.ListenAndServe(s.IpAddr, r)
}

func (s *Server) SignUpHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if r.Method != http.MethodPost {
		return fmt.Errorf("invalid mathod")
	}

	if r.Body == nil {
		return fmt.Errorf("request is null")
	}
	defer r.Body.Close()

	user := &userservice.UserModel{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return err
	}

	resp, err := s.userServiceClient.SignUp(ctx, user.Name, user.Password, user.Email, user.Phone, user.UserType)
	if err != nil {
		return err
	}

	return WriteJson(w, http.StatusOK, resp)
}

func (s *Server) SignInHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if r.Method != http.MethodGet {
		return fmt.Errorf("invalid mathod")
	}

	if r.Body == nil {
		return fmt.Errorf("request is null")
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

	claims, err := ValidateToken(r_token)
	if err != nil {
		return err
	}

	newToken, err := s.userServiceClient.GetAccessToken(ctx, claims.User_id, r_token)
	if err != nil {
		return err
	}
	return WriteJson(w, http.StatusOK, newToken)
}

type GetHandlerFunc func(w http.ResponseWriter, r *http.Request) error

func createHandlerFunc(f GetHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			log.Println("Handler error", err)
			return
		}
	}
}

func WriteJson(w http.ResponseWriter, status int, msg any) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(msg)
}
