package books

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	elastic "github.com/elastic/go-elasticsearch"
)

type Repository interface {
	CreateBook(ctx context.Context, book Book) error
	UpdateBook(ctx context.Context, book Book) error
	DeleteBook(ctx context.Context, isbn string) error

	GetBookByISBN(ctx context.Context, isbn string) (*Book, error)
	GetBooks(ctx context.Context, limit int, offset int) (*BookList, error)
	SearchBook(ctx context.Context, query string, limit int, offset int) (*BookList, error)
}

type wrapper struct {
	Source Book `json:"_source"`
}

type booksRepository struct {
	client *elastic.Client
}

func NewRepository(dsn string) (Repository, error) {
	client, err := elastic.NewClient(elastic.Config{
		Addresses: []string{
			dsn,
		},
	})

	if err != nil {
		LogError(slog.LevelError, "Database", err, "New Repository Error")
		return nil, err
	}
	LogInfo(slog.LevelInfo, "Database", fmt.Sprintf("Database Created at url %s", dsn))
	return &booksRepository{
		client: client,
	}, nil
}

func (r *booksRepository) CreateBook(ctx context.Context, book Book) error {
	var jsonIndex bytes.Buffer
	err := json.NewEncoder(&jsonIndex).Encode(book)
	if err != nil {
		LogError(slog.LevelError, "Database", fmt.Errorf("encoding error %s", err), "Create Book Error")
		return err
	}
	_, err = r.client.Index(
		"books",
		&jsonIndex,
		r.client.Index.WithDocumentID(book.ISBN),
		r.client.Index.WithContext(ctx),
	)

	if err != nil {
		LogError(slog.LevelError, "Database", err, "Create Book Error")
		return err
	}

	LogInfo(slog.LevelInfo, "Database", "Book Created")
	return nil
}

func (r *booksRepository) UpdateBook(ctx context.Context, book Book) error {
	err := r.DeleteBook(ctx, book.ISBN)
	if err != nil {
		return err
	}

	err = r.CreateBook(ctx, book)
	if err != nil {
		return err
	}
	LogInfo(slog.LevelInfo, "Database", "Book Updated")
	return nil
}

func (r *booksRepository) DeleteBook(ctx context.Context, isbn string) error {
	_, err := r.client.Delete(
		"books",
		isbn,
		r.client.Delete.WithContext(ctx),
	)
	if err != nil {
		LogError(slog.LevelError, "Database", err, "Delete Book Error")
		return err
	}
	LogInfo(slog.LevelInfo, "Database", "Book deleted")
	return nil
}

func (r *booksRepository) GetBookByISBN(ctx context.Context, isbn string) (*Book, error) {
	res, err := r.client.Get(
		"books",
		isbn,
		r.client.Get.WithContext(ctx),
	)
	if err != nil {
		LogError(slog.LevelError, "Database", err, "Get Book by ISBN Error")
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		LogError(slog.LevelError, "Database", fmt.Errorf("entity not found %s", err), "Get Book by ISBN Error")
		return nil, err
	}

	wr := &wrapper{}
	err = json.NewDecoder(res.Body).Decode(&wr)
	if err != nil {
		LogError(slog.LevelError, "Database", fmt.Errorf("decoding error %s", err), "Get Book by ISBN Error")
		return nil, err
	}
	LogInfo(slog.LevelInfo, "Database", "Got Book By ISBN")
	return &wr.Source, nil
}

func (r *booksRepository) GetBooks(ctx context.Context, limit int, offset int) (*BookList, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
		"from": offset,
		"size": limit,
	}
	var queryBuf bytes.Buffer
	err := json.NewEncoder(&queryBuf).Encode(query)
	if err != nil {
		LogError(slog.LevelError, "Database", fmt.Errorf("encoding error %s", err), "Get Book Error")
		return nil, err
	}

	res, err := r.client.Search(
		r.client.Search.WithIndex("books"),
		r.client.Search.WithBody(&queryBuf),
		r.client.Search.WithContext(ctx),
	)
	if err != nil {
		LogError(slog.LevelError, "Database", err, "Get All Books Error")
		return nil, err
	}

	if res.IsError() {
		LogError(slog.LevelError, "Database", fmt.Errorf("entity not found %s", err), "Get All Books Error")
		return nil, err
	}
	defer res.Body.Close()

	var esRes struct {
		Hits struct {
			Total struct {
				Value int `json:"value"` // Total number of hits
			} `json:"total"`
			Hits []struct {
				Source Book `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&esRes); err != nil {
		LogError(slog.LevelError, "Database", fmt.Errorf("decoding error %s", err), "Delete All Books Error")
		return nil, err
	}

	books := make([]Book, 0, len(esRes.Hits.Hits))
	for _, v := range esRes.Hits.Hits {
		books = append(books, v.Source)
	}

	totalBook := esRes.Hits.Total.Value
	totalPage := totalBook / limit
	if totalBook%limit != 0 {
		totalPage++
	}

	LogInfo(slog.LevelInfo, "Database", "Got All Books")
	LogInfo(slog.LevelInfo, "Database", BookList{
		TotalPage:  totalPage,
		TotalBooks: totalBook,
		Books:      books,
	})

	return &BookList{
		TotalPage:  totalPage,
		TotalBooks: totalBook,
		Books:      books,
	}, nil
}

func (r *booksRepository) SearchBook(ctx context.Context, query string, limit int, offset int) (*BookList, error) {
	q := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"title": map[string]interface{}{
					"query":     query,
					"fuzziness": "AUTO", // Optional: allows for typo-tolerant search
				},
			},
		},
		"from": offset,
		"size": limit,
	}

	var queryBuf bytes.Buffer
	err := json.NewEncoder(&queryBuf).Encode(q)
	if err != nil {
		LogError(slog.LevelError, "Database", fmt.Errorf("encoding error %s", err), "Search Book Error")
		return nil, err
	}

	res, err := r.client.Search(
		r.client.Search.WithIndex("books"),
		r.client.Search.WithBody(&queryBuf),
		r.client.Search.WithContext(ctx),
		r.client.Search.WithTrackTotalHits(true),
		r.client.Search.WithPretty(),
	)

	if err != nil {
		LogError(slog.LevelError, "Database", err, "Search Books Error")
		return nil, err
	}

	if res.IsError() {
		LogError(slog.LevelError, "Database", fmt.Errorf("entity not found %s", err), "Search Books Error")
		return nil, err
	}
	defer res.Body.Close()

	var esRes struct {
		Hits struct {
			Total struct {
				Value int `json:"value"` // total number of hits
			} `json:"total"`
			Hits []struct {
				Source Book `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	err = json.NewDecoder(res.Body).Decode(&esRes)
	if err != nil {
		LogError(slog.LevelError, "Database", err, "Decoding failed")
		return nil, err
	}

	books := make([]Book, 0, len(esRes.Hits.Hits))
	for _, v := range esRes.Hits.Hits {
		books = append(books, v.Source)
	}

	totalBook := esRes.Hits.Total.Value
	totalPage := totalBook / limit
	if totalBook%limit != 0 {
		totalPage++
	}

	LogInfo(slog.LevelInfo, "Database", fmt.Sprintf("Got All Books By Query=%s", query))
	return &BookList{
		TotalPage:  totalPage,
		TotalBooks: totalBook,
		Books:      books,
	}, nil
}
