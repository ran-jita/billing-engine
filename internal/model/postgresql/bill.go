package postgresql

import "time"

type Bill struct {
	ID         string    `json:"id" db:"id"`
	BorrowerId string    `json:"borrower_id" db:"borrower_id"`
	LoanId     string    `json:"loan_id" db:"loan_id"`
	Amount     float64   `json:"amount" db:"amount"`
	DueDate    time.Time `json:"due_date" db:"due_date"`
	Status     string    `json:"status" db:"status"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}
