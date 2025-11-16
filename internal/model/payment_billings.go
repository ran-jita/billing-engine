package model

import "time"

type PaymentBilling struct {
	ID        string    `json:"id" db:"id"`
	PaymentId string    `json:"payment_id" db:"payment_id"`
	BillingId string    `json:"billing_id" db:"billing_id"`
	Amount    float64   `json:"amount" db:"amount"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
