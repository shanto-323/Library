package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/shanto-323/Library/books"
)

func (s *Server) CreateBookHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if r.Method != http.MethodPost {
		return fmt.Errorf("invalid mathod")
	}

	u_type := r.Context().Value(userTypeKey).(string)
	if u_type != "ADMIN" {
		return fmt.Errorf("not allowed")
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

	return WriteJson(w, http.StatusOK, b)
}

func (s *Server) UpdateBookHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if r.Method != http.MethodPatch {
		return fmt.Errorf("invalid mathod")
	}

	u_type := r.Context().Value(userTypeKey).(string)
	if u_type != "ADMIN" {
		return fmt.Errorf("not allowed")
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
	return WriteJson(w, http.StatusOK, &msg)
}

func (s *Server) DeleteBookHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	if r.Method != http.MethodDelete {
		return fmt.Errorf("invalid mathod")
	}

	u_type := r.Context().Value(userTypeKey).(string)
	if u_type != "ADMIN" {
		return fmt.Errorf("not allowed")
	}

	isbn := mux.Vars(r)["isbn"]
	if isbn == "" {
		return fmt.Errorf("missing ISBN in URL")
	}

	msg, err := s.bookClient.DeleteBook(ctx, isbn)
	if err != nil {
		return err
	}

	return WriteJson(w, http.StatusOK, &msg)
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

	return WriteJson(w, http.StatusOK, b)
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
	return WriteJson(w, http.StatusOK, books)
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
	return WriteJson(w, http.StatusOK, books)
}
