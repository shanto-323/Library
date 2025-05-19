package books

import (
	"context"
	"fmt"
	"time"
)

type Service interface {
	NewBook(ctx context.Context, title string, isbn string, writer string, t_copies uint64, on_loan uint64) (*Book, error)
	UpdateBook(ctx context.Context, book Book) error
	DeleteBook(ctx context.Context, isbn string) error
	GetBook(ctx context.Context, isbn string) (*Book, error)
	GetAllBook(ctx context.Context, limit uint64, offset uint64) ([]Book, error)
	SearchBook(ctx context.Context, query string, limit uint64, offset uint64) ([]Book, error)
}

type bookService struct {
	repository Repository
}

func NewBookService(repository Repository) Service {
	return &bookService{repository: repository}
}

func (s *bookService) NewBook(ctx context.Context, title string, isbn string, writer string, t_copies uint64, on_loan uint64) (*Book, error) {
	newBook := Book{
		ISBN:        isbn,
		Title:       title,
		Writer:      writer,
		TotalCopies: t_copies,
		OnLoan:      on_loan,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := s.repository.CreateBook(ctx, newBook)
	if err != nil {
		return nil, fmt.Errorf("bookService error %s", err)
	}
	return &newBook, nil
}

func (s *bookService) UpdateBook(ctx context.Context, book Book) error {
	err := s.repository.UpdateBook(ctx, book)
	if err != nil {
		return err
	}
	return nil
}

func (s *bookService) DeleteBook(ctx context.Context, isbn string) error {
	err := s.repository.DeleteBook(ctx, isbn)
	if err != nil {
		return err
	}
	return nil
}

func (s *bookService) GetBook(ctx context.Context, isbn string) (*Book, error) {
	book, err := s.repository.GetBookByISBN(ctx, isbn)
	if err != nil {
		return nil, fmt.Errorf("bookService error %s", err)
	}
	return book, nil
}

func (s *bookService) GetAllBook(ctx context.Context, limit uint64, offset uint64) ([]Book, error) {
	if limit < 10 || limit > 100 {
		limit = 10
	}
	books, err := s.repository.GetBooks(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	return books, nil
}

func (s *bookService) SearchBook(ctx context.Context, query string, limit uint64, offset uint64) ([]Book, error) {
	if limit < 10 || limit > 100 {
		limit = 10
	}
	books, err := s.repository.SearchBook(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	return books, nil
}
