package postgresql

import "time"

type Borrower struct {
	ID         string    `json:"id" db:"id"`
	Name       string    `json:"name" db:"name"`
	Delinquent bool      `json:"delinquent" db:"delinquent"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}
