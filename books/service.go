package books

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

type Service interface {
	NewBook(ctx context.Context, title string, isbn string, writer string, t_copies uint64, on_loan uint64) (*Book, error)
	UpdateBook(ctx context.Context, book *Book) error
	DeleteBook(ctx context.Context, isbn string) error
	GetBook(ctx context.Context, isbn string) (*Book, error)
	GetAllBook(ctx context.Context, limit int64, offset int64) (*BookList, error)
	SearchBook(ctx context.Context, query string, limit int64, offset int64) (*BookList, error)
}

type bookService struct {
	repository Repository
}

func NewBookService(repository Repository) Service {
	return &bookService{repository: repository}
}

func (s *bookService) NewBook(ctx context.Context, title string, isbn string, writer string, t_copies uint64, on_loan uint64) (*Book, error) {
	if book, _ := s.GetBook(ctx, isbn); book != nil {
		return nil, fmt.Errorf("book exists")
	}

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

func (s *bookService) UpdateBook(ctx context.Context, book *Book) error {
	dbRecord, err := s.GetBook(ctx, book.ISBN)
	if err != nil {
		return err
	}

	updatedBook, hasChange := mutationHelper(dbRecord, book)
	if !hasChange {
		return fmt.Errorf("no new changes")
	}

	err = s.repository.UpdateBook(ctx, *updatedBook)
	if err != nil {
		return err
	}
	LogInfo(slog.LevelInfo, "Service", "Book Created")
	LogInfo(slog.LevelInfo, "Service", updatedBook)
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

func (s *bookService) GetAllBook(ctx context.Context, limit int64, offset int64) (*BookList, error) {
	if limit < 10 || limit > 100 {
		limit = 10
	}
	books, err := s.repository.GetBooks(ctx, int(limit), int(offset))
	if err != nil {
		return nil, err
	}
	return books, nil
}

func (s *bookService) SearchBook(ctx context.Context, query string, limit int64, offset int64) (*BookList, error) {
	if limit < 10 || limit > 100 {
		limit = 10
	}
	books, err := s.repository.SearchBook(ctx, query, int(limit), int(offset))
	if err != nil {
		return nil, err
	}
	return books, nil
}

func mutationHelper(dbRecord *Book, book *Book) (*Book, bool) {
	hasChange := false
	if book.Title != "" && book.Title != dbRecord.Title {
		dbRecord.Title = book.Title
		hasChange = true
	}
	if book.Writer != "" && book.Writer != dbRecord.Writer {
		dbRecord.Writer = book.Writer
		hasChange = true
	}
	if book.TotalCopies != dbRecord.TotalCopies {
		dbRecord.TotalCopies = book.TotalCopies
		hasChange = true
	}
	if book.OnLoan != dbRecord.OnLoan {
		dbRecord.OnLoan = book.OnLoan
		hasChange = true
	}

	dbRecord.UpdatedAt = time.Now()
	return dbRecord, hasChange
}
