package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/shanto-323/Library/books"
)

type Server struct {
	IpAddr     string
	bookClient *books.Client
}

func NewServer(ipAddr string, bookClientUrl string) (*Server, error) {
	bookClient, err := books.NewClient(bookClientUrl)
	if err != nil {
		log.Println("NewServer Error", err)
		return nil, err
	}

	return &Server{
		IpAddr:     ipAddr,
		bookClient: bookClient,
	}, nil
}

func (s *Server) Start() error {
	r := mux.NewRouter()
	router := r.PathPrefix("/library/v1").Subrouter()

	//books
	bookRouter := router.PathPrefix("/books").Subrouter()

	bookRouter.HandleFunc("/search", createHandlerFunc(s.SearchBookHandler)).Methods("GET")
	bookRouter.HandleFunc("", createHandlerFunc(s.CreateBookHandler)).Methods("POST")
	bookRouter.HandleFunc("/{isbn}", createHandlerFunc(s.UpdateBookHandler)).Methods("PATCH")
	bookRouter.HandleFunc("/{isbn}", createHandlerFunc(s.DeleteBookHandler)).Methods("DELETE")
	bookRouter.HandleFunc("/{isbn}", createHandlerFunc(s.GetBookHandler)).Methods("GET")
	bookRouter.HandleFunc("", createHandlerFunc(s.GetAllBooksHandler)).Methods("GET")

	fmt.Println("Api running.. ")
	return http.ListenAndServe(s.IpAddr, r)
}

func (s *Server) CreateBookHandler(w http.ResponseWriter, r *http.Request) error {
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

	book := &books.Book{}
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		return err
	}

	b, err := s.bookClient.CreateBook(ctx, book.Title, book.ISBN, book.Writer, book.TotalCopies, book.OnLoan)
	if err != nil {
		return err
	}

	return WriteJson(w, b)
}

func (s *Server) UpdateBookHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	if r.Method != http.MethodPatch {
		return fmt.Errorf("invalid mathod")
	}

	if r.Body == nil {
		return fmt.Errorf("request is null")
	}
	defer r.Body.Close()

	book := &books.Book{}
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		return err
	}

	book.ISBN = mux.Vars(r)["isbn"]
	if book.ISBN == "" {
		return fmt.Errorf("missing ISBN in URL")
	}

	msg, err := s.bookClient.UpdateBook(ctx, book)
	if err != nil {
		return err
	}

	books.LogInfo(slog.LevelInfo, "MAIN SERVER", book)
	return WriteJson(w, &msg)
}

func (s *Server) DeleteBookHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	if r.Method != http.MethodDelete {
		return fmt.Errorf("invalid mathod")
	}

	isbn := mux.Vars(r)["isbn"]
	if isbn == "" {
		return fmt.Errorf("missing ISBN in URL")
	}

	msg, err := s.bookClient.DeleteBook(ctx, isbn)
	if err != nil {
		return err
	}

	return WriteJson(w, &msg)
}

func (s *Server) GetBookHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if r.Method != http.MethodGet {
		return fmt.Errorf("invalid mathod")
	}

	isbn := mux.Vars(r)["isbn"]
	if isbn == "" {
		return fmt.Errorf("missing ISBN in URL")
	}

	b, err := s.bookClient.GetBook(ctx, isbn)
	if err != nil {
		return err
	}

	return WriteJson(w, b)
}

func (s *Server) GetAllBooksHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if r.Method != http.MethodGet {
		return fmt.Errorf("invalid mathod")
	}

	limit, err := perseInt("limit", r)
	if err != nil {
		return err
	}
	offset, err := perseInt("offset", r)
	if err != nil {
		return err
	}

	books, err := s.bookClient.GetAllBook(ctx, limit, offset)
	if err != nil {
		return err
	}
	return WriteJson(w, books)
}

func (s *Server) SearchBookHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if r.Method != http.MethodGet {
		return fmt.Errorf("invalid mathod")
	}

	limit, err := perseInt("limit", r)
	if err != nil {
		return err
	}
	offset, err := perseInt("offset", r)
	if err != nil {
		return err
	}
	query := r.URL.Query().Get("query")
	if query == "" {
		return nil
	}

	books, err := s.bookClient.SearchBook(ctx, query, limit, offset)
	if err != nil {
		return err
	}
	return WriteJson(w, books)
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

func WriteJson(w http.ResponseWriter, msg any) error {
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
