package books

import (
	"context"
	"fmt"

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
		return nil, err
	}

	return &Book{
		Title:       book.Book.Title,
		ISBN:        book.Book.Isbn,
		Writer:      book.Book.Writer,
		TotalCopies: book.Book.TotalCopies,
		OnLoan:      book.Book.OnLoan,
		CreatedAt:   book.Book.CreatedAt.AsTime(),
		UpdatedAt:   book.Book.UpdatedAt.AsTime(),
	}, nil
}

func (c *Client) UpdateBook(ctx context.Context, book *Book) (*string, error) {
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
		return nil, err
	}
	return &resp.Msg, nil
}

func (c *Client) DeleteBook(ctx context.Context, isbn string) (*string, error) {
	resp, err := c.service.DeleteBook(
		ctx,
		&pb.DeleteBookRequest{
			Isbn: isbn,
		},
	)
	if err != nil {
		return nil, err
	}
	return &resp.Msg, nil
}

func (c *Client) GetBook(ctx context.Context, isbn string) (*Book, error) {
	book, err := c.service.GetBook(ctx, &pb.GetBookRequest{Isbn: isbn})
	if err != nil {
		return nil, err
	}

	return &Book{
		Title:       book.Book.Title,
		ISBN:        book.Book.Isbn,
		Writer:      book.Book.Writer,
		TotalCopies: book.Book.TotalCopies,
		OnLoan:      book.Book.OnLoan,
		CreatedAt:   book.Book.CreatedAt.AsTime(),
		UpdatedAt:   book.Book.UpdatedAt.AsTime(),
	}, nil
}

func (c *Client) GetAllBook(ctx context.Context, limit int64, offset int64) (*BookList, error) {
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
	for _, b := range resp.Books {
		book := Book{
			Title:       b.Title,
			ISBN:        b.Isbn,
			Writer:      b.Writer,
			TotalCopies: b.TotalCopies,
			OnLoan:      b.OnLoan,
			CreatedAt:   b.CreatedAt.AsTime(),
			UpdatedAt:   b.UpdatedAt.AsTime(),
		}
		books = append(books, book)
	}

	return &BookList{
		TotalPage:  int(resp.TotalPages),
		TotalBooks: int(resp.TotalBooks),
		Books:      books,
	}, nil
}
func (c *Client) SearchBook(ctx context.Context, query string, limit int64, offset int64) (*BookList, error) {
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
	for _, b := range resp.Books {
		book := Book{
			Title:       b.Title,
			ISBN:        b.Isbn,
			Writer:      b.Writer,
			TotalCopies: b.TotalCopies,
			OnLoan:      b.OnLoan,
			CreatedAt:   b.CreatedAt.AsTime(),
			UpdatedAt:   b.UpdatedAt.AsTime(),
		}
		books = append(books, book)
	}

	return &BookList{
		TotalPage:  int(resp.TotalPages),
		TotalBooks: int(resp.TotalBooks),
		Books:      books,
	}, nil
}
