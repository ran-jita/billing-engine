package dto

type CreatePayment struct {
	TotalAmount float64 `json:"total_amount" db:"total_amount"`
	LoandId     string  `json:"loand_id" db:"loand_id"`
}
