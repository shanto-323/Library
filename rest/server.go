package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

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

	// books
	bookRouter := router.PathPrefix("/books").Subrouter()
	bookRouter.Use(JwtMiddleWere)

	bookRouter.HandleFunc("/search", createHandlerFunc(s.SearchBookHandler)).Methods("GET")
	bookRouter.HandleFunc("/{isbn}", createHandlerFunc(s.GetBookHandler)).Methods("GET")
	bookRouter.HandleFunc("", createHandlerFunc(s.GetAllBooksHandler)).Methods("GET")

	bookAdminRouter := bookRouter.PathPrefix("/admin").Subrouter()
	bookAdminRouter.Use(JwtMiddleWere)
	bookAdminRouter.HandleFunc("", createHandlerFunc(s.CreateBookHandler)).Methods("POST")
	bookAdminRouter.HandleFunc("/{isbn}", createHandlerFunc(s.UpdateBookHandler)).Methods("PATCH")
	bookAdminRouter.HandleFunc("/{isbn}", createHandlerFunc(s.DeleteBookHandler)).Methods("DELETE")

	// user
	userServiceRouter := router.PathPrefix("/user").Subrouter()
	userServiceRouter.HandleFunc("/signup", createHandlerFunc(s.SignUpHandler)).Methods("POST")
	userServiceRouter.HandleFunc("/login", createHandlerFunc(s.SignInHandler)).Methods("POST")
	userServiceRouter.HandleFunc("/token/{id}", createHandlerFunc(s.NewRefreshTokenHandler)).Methods("GET")

	authUserServiceRouter := userServiceRouter.PathPrefix("").Subrouter()
	authUserServiceRouter.Use(JwtMiddleWere)
	authUserServiceRouter.HandleFunc("/logout/{id}", createHandlerFunc(s.SignOutHandler)).Methods("POST")
	authUserServiceRouter.HandleFunc("/{id}", createHandlerFunc(s.UpdateUserHandler)).Methods("PATCH")
	authUserServiceRouter.HandleFunc("/delete/{id}", createHandlerFunc(s.DeleteUserHandler)).Methods("DELETE")

	// user -ADMIN Gateway
	protectedUserServiceRouter := userServiceRouter.PathPrefix("/admin").Subrouter()
	protectedUserServiceRouter.Use(JwtMiddleWere)
	protectedUserServiceRouter.HandleFunc("/all", createHandlerFunc(s.GetALlUserHandler)).Methods("GET")

	fmt.Println("Api running.. ")
	return http.ListenAndServe(s.IpAddr, r)
}

type GetHandlerFunc func(w http.ResponseWriter, r *http.Request) error

func createHandlerFunc(f GetHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJson(w, http.StatusBadRequest, err)
			return
		}
	}
}

func WriteJson(w http.ResponseWriter, status int, msg any) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(msg)
}

func perseInt(v string, r *http.Request) (int64, error) {
	qv := r.URL.Query().Get(v)
	if qv == "" {
		return 0, nil
	}
	num, err := strconv.ParseInt(qv, 10, 64)
	if err != nil {
		return 0, nil
	}

	return num, nil
}
