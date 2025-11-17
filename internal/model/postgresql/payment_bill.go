package postgresql

import "time"

type PaymentBill struct {
	ID        string    `json:"id" db:"id"`
	PaymentId string    `json:"payment_id" db:"payment_id"`
	BillId    string    `json:"bill_id" db:"bill_id"`
	Amount    float64   `json:"amount" db:"amount"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
