package dto

type CreatePayment struct {
	TotalAmount float64 `json:"total_amount"`
	LoanId      string  `json:"loan_id"`
	PaymentDate string  `json:"payment_date"`
}
