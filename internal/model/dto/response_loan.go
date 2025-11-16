package dto

import "github.com/ran-jita/billing-engine/internal/model"

type LoanWithBillings struct {
	Loan     model.Loan      `json:"loan"`
	Billings []model.Billing `json:"billings"`
}
