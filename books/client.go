package books

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/shanto-323/Library/books/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.BookServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("client error %s", err)
	}
	s := pb.NewBookServiceClient(conn)
	return &Client{conn: conn, service: s}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) CreateBook(ctx context.Context, title string, isbn string, writer string, t_copies uint64, on_loan uint64) (*Book, error) {
	book, err := c.service.CreateBook(
		ctx,
		&pb.CreateBookRequest{
			Title:       title,
			Isbn:        isbn,
			Writer:      writer,
			TotalCopies: t_copies,
			OnLoan:      on_loan,
		},
	)
	if err != nil {
		SlogLogger(
			LogType{
				L: slog.LevelError,
				M: "Client",
				E: err,
				D: "Create book",
			},
		)
		return nil, err
	}
	c_at, _ := time.Parse(time.RFC3339, book.Book.CratedAt)
	u_at, err := time.Parse(time.RFC3339, book.Book.UpdatedAt)
	if err != nil {
		SlogLogger(
			LogType{
				L: slog.LevelError,
				M: "Client",
				E: fmt.Errorf("time persing error %s", err),
				D: "Create book",
			},
		)
	}

	return &Book{
		Title:       book.Book.Title,
		ISBN:        book.Book.Isbn,
		Writer:      book.Book.Writer,
		TotalCopies: book.Book.TotalCopies,
		OnLoan:      book.Book.OnLoan,
		CreatedAt:   c_at,
		UpdatedAt:   u_at,
	}, nil
}

func (c *Client) UpdateBook(ctx context.Context, book Book) error {
	resp, err := c.service.UpdateBook(
		ctx,
		&pb.UpdateBookRequest{
			Book: &pb.Book{
				Title:       book.Title,
				Isbn:        book.ISBN,
				Writer:      book.Writer,
				TotalCopies: book.TotalCopies,
				OnLoan:      book.OnLoan,
			},
		},
	)

	if err != nil {
		return err
	}
	SlogLogger(
		LogType{
			L: slog.LevelInfo,
			M: "Client",
			D: resp.Msg,
		},
	)
	return nil
}

func (c *Client) DeleteBook(ctx context.Context, isbn string) error {
	resp, err := c.service.DeleteBook(
		ctx,
		&pb.DeleteBookRequest{
			Isbn: isbn,
		},
	)

	if err != nil {
		return err
	}
	SlogLogger(
		LogType{
			L: slog.LevelInfo,
			M: "Client",
			D: resp.Msg,
		},
	)
	return nil
}

func (c *Client) GetBook(ctx context.Context, isbn string) (*Book, error) {
	book, err := c.service.GetBook(ctx, &pb.GetBookRequest{Isbn: isbn})
	if err != nil {
		return nil, err
	}

	c_at, _ := time.Parse(time.RFC3339, book.Book.CratedAt)
	u_at, err := time.Parse(time.RFC3339, book.Book.UpdatedAt)
	if err != nil {
		SlogLogger(
			LogType{
				L: slog.LevelError,
				M: "Client",
				E: fmt.Errorf("time persing error %s", err),
				D: "Create book",
			},
		)
	}

	return &Book{
		Title:       book.Book.Title,
		ISBN:        book.Book.Isbn,
		Writer:      book.Book.Writer,
		TotalCopies: book.Book.TotalCopies,
		OnLoan:      book.Book.OnLoan,
		CreatedAt:   c_at,
		UpdatedAt:   u_at,
	}, nil
}

func (c *Client) GetAllBook(ctx context.Context, limit uint64, offset uint64) ([]Book, error) {
	resp, err := c.service.GetAllBook(
		ctx,
		&pb.GetAllBookRequest{
			Limit:  limit,
			Offset: offset,
		},
	)

	if err != nil {
		return nil, err
	}

	books := []Book{}
	for _, b := range resp.Book {
		c_at, _ := time.Parse(time.RFC3339, b.CratedAt)
		u_at, err := time.Parse(time.RFC3339, b.UpdatedAt)
		if err != nil {
			SlogLogger(
				LogType{
					L: slog.LevelError,
					M: "Client",
					E: fmt.Errorf("time persing error %s", err),
					D: "Get All Books",
				},
			)
		}
		book := Book{
			Title:       b.Title,
			ISBN:        b.Isbn,
			Writer:      b.Writer,
			TotalCopies: b.TotalCopies,
			OnLoan:      b.OnLoan,
			CreatedAt:   c_at,
			UpdatedAt:   u_at,
		}
		books = append(books, book)
	}

	return books, nil
}
func (c *Client) SearchBook(ctx context.Context, query string, limit uint64, offset uint64) ([]Book, error) {
	resp, err := c.service.SearchBook(
		ctx,
		&pb.SearchBookRequest{
			Query:  query,
			Limit:  limit,
			Offset: offset,
		},
	)

	if err != nil {
		return nil, err
	}

	books := []Book{}
	for _, b := range resp.Book {
		c_at, _ := time.Parse(time.RFC3339, b.CratedAt)
		u_at, err := time.Parse(time.RFC3339, b.UpdatedAt)
		if err != nil {
			SlogLogger(
				LogType{
					L: slog.LevelError,
					M: "Client",
					E: fmt.Errorf("time persing error %s", err),
					D: "Get All Books",
				},
			)
		}
		book := Book{
			Title:       b.Title,
			ISBN:        b.Isbn,
			Writer:      b.Writer,
			TotalCopies: b.TotalCopies,
			OnLoan:      b.OnLoan,
			CreatedAt:   c_at,
			UpdatedAt:   u_at,
		}
		books = append(books, book)
	}

	return books, nil
}
