package books

import "time"

type Book struct {
	Title       string    `json:"title"`
	ISBN        string    `json:"isbn"`
	Writer      string    `json:"writer"`
	TotalCopies uint64    `json:"total_copies"`
	OnLoan      uint64    `json:"on_loan"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
