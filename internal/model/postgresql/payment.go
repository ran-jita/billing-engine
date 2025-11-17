package postgresql

import "time"

type Payment struct {
	ID          string    `json:"id" db:"id"`
	BorrowerID  string    `json:"borrower_id" db:"borrower_id"`
	TotalAmount float64   `json:"total_amount" db:"total_amount"`
	PaymentDate time.Time `json:"delinquent" db:"delinquent"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
