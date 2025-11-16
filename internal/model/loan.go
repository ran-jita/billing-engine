package model

import "time"

type Loan struct {
	ID                string    `json:"id" db:"id"`
	BorrowerID        string    `json:"borrower_id" db:"borrower_id"`
	BaseAmount        float64   `json:"base_amount" db:"base_amount"`
	FeeAmount         float64   `json:"fee_amount" db:"fee_amount"`
	TotalLoanAmount   float64   `json:"total_loan_amount" db:"total_loan_amount"`
	OutstandingAmount float64   `json:"outstanding_amount" db:"outstanding_amount"`
	StartDate         time.Time `json:"start_date" db:"start_date"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
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
