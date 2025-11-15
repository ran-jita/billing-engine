package domain

import "time"

type Loan struct {
	ID              string    `json:"id"`
	BorrowerID      string    `json:"borrower_id"`
	PrincipalAmount float64   `json:"principal_amount"`
	InterestRate    float64   `json:"interest_rate"`
	Term            int       `json:"term"` // dalam bulan
	StartDate       time.Time `json:"start_date"`
	Status          string    `json:"status"` // active, paid, defaulted
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type Billing struct {
	ID              string    `json:"id"`
	LoanID          string    `json:"loan_id"`
	DueDate         time.Time `json:"due_date"`
	Amount          float64   `json:"amount"`
	PrincipalAmount float64   `json:"principal_amount"`
	InterestAmount  float64   `json:"interest_amount"`
	Status          string    `json:"status"` // pending, paid, overdue
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type Payment struct {
	ID          string    `json:"id"`
	BillingID   string    `json:"billing_id"`
	Amount      float64   `json:"amount"`
	PaymentDate time.Time `json:"payment_date"`
	Method      string    `json:"method"` // transfer, cash, etc
	Status      string    `json:"status"` // success, pending, failed
	CreatedAt   time.Time `json:"created_at"`
}
