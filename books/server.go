package books

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/shanto-323/Library/books/pb"
	"google.golang.org/grpc"
)

type grpcServer struct {
	pb.UnimplementedBookServiceServer
	service Service
}

func ListenGRPC(s Service, port string) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {

	}
	serv := grpc.NewServer()
	pb.RegisterBookServiceServer(serv, &grpcServer{
		service: s,
	})

	return serv.Serve(lis)
}

func (sr *grpcServer) CreateBook(ctx context.Context, r *pb.CreateBookRequest) (*pb.CreateBookResponse, error) {
	book, err := sr.service.NewBook(ctx, r.Title, r.Isbn, r.Writer, r.TotalCopies, r.OnLoan)
	if err != nil {
		SlogLogger(
			LogType{
				L: slog.LevelError,
				M: "Server",
				E: err,
				D: "New book",
			},
		)
		return nil, err
	}
	return &pb.CreateBookResponse{Book: &pb.Book{
		Title:       book.Title,
		Isbn:        book.ISBN,
		Writer:      book.Writer,
		TotalCopies: book.TotalCopies,
		OnLoan:      book.OnLoan,
		CratedAt:    book.CreatedAt.Format(time.RFC3339), //use time  book.Book
		UpdatedAt:   book.UpdatedAt.Format(time.RFC3339),
	}}, nil
}

func (sr *grpcServer) UpdateBook(ctx context.Context, r *pb.UpdateBookRequest) (*pb.UpdateBookResponse, error) {
	book := r.Book
	err := sr.service.UpdateBook(ctx, Book{
		Title:       book.Title,
		ISBN:        book.Isbn,
		Writer:      book.Writer,
		TotalCopies: book.TotalCopies,
		OnLoan:      book.OnLoan,
		UpdatedAt:   time.Now(),
	})

	if err != nil {
		return nil, err
	}
	return &pb.UpdateBookResponse{
		Msg: "book updated succesfully",
	}, nil
}

func (sr *grpcServer) DeleteBook(ctx context.Context, r *pb.DeleteBookRequest) (*pb.DeleteBookResponse, error) {
	err := sr.service.DeleteBook(ctx, r.Isbn)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteBookResponse{
		Msg: "book deleted succesfully",
	}, nil
}

func (sr *grpcServer) GetBook(ctx context.Context, r *pb.GetBookRequest) (*pb.GetBookResponse, error) {
	if r.Isbn == "" {
		SlogLogger(
			LogType{
				L: slog.LevelError,
				M: "Server",
				E: fmt.Errorf("isbn is empty"),
				D: "Get Book",
			},
		)
	}

	book, err := sr.service.GetBook(ctx, r.Isbn)
	if err != nil {
		return nil, err
	}

	if book == nil {
		if r.Isbn == "" {
			SlogLogger(
				LogType{
					L: slog.LevelError,
					M: "Server",
					E: fmt.Errorf("isbn is empty %s", r.Isbn),
					D: "Get Book",
				},
			)
		}
		return nil, fmt.Errorf("isbn is empty %s", r.Isbn)
	}

	return &pb.GetBookResponse{Book: &pb.Book{
		Title:       book.Title,
		Isbn:        book.ISBN,
		Writer:      book.Writer,
		TotalCopies: book.TotalCopies,
		OnLoan:      book.OnLoan,
		CratedAt:    book.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   book.UpdatedAt.Format(time.RFC3339),
	}}, nil
}

func (sr *grpcServer) GetAllBook(ctx context.Context, r *pb.GetAllBookRequest) (*pb.GetAllBookResponse, error) {
	books, err := sr.service.GetAllBook(ctx, r.Limit, r.Offset)
	if err != nil {
		return nil, err
	}

	pbBook := []*pb.Book{}
	for _, book := range books {
		b := &pb.Book{
			Title:       book.Title,
			Isbn:        book.ISBN,
			Writer:      book.Writer,
			TotalCopies: book.TotalCopies,
			OnLoan:      book.OnLoan,
			CratedAt:    book.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   book.UpdatedAt.Format(time.RFC3339),
		}
		pbBook = append(pbBook, b)
	}

	return &pb.GetAllBookResponse{
		Book: pbBook,
	}, nil
}

func (sr *grpcServer) SearchBook(ctx context.Context, r *pb.SearchBookRequest) (*pb.SearchBookResponse, error) {
	books, err := sr.service.SearchBook(ctx, r.Query, r.Limit, r.Offset)
	if err != nil {
		return nil, err
	}

	pbBook := []*pb.Book{}
	for _, book := range books {
		b := &pb.Book{
			Title:       book.Title,
			Isbn:        book.ISBN,
			Writer:      book.Writer,
			TotalCopies: book.TotalCopies,
			OnLoan:      book.OnLoan,
			CratedAt:    book.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   book.UpdatedAt.Format(time.RFC3339),
		}
		pbBook = append(pbBook, b)
	}

	return &pb.SearchBookResponse{
		Book: pbBook,
	}, nil
}
