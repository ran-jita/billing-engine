package dto

import "github.com/ran-jita/billing-engine/internal/model"

type LoanWithBills struct {
	Loan  model.Loan   `json:"loan"`
	Bills []model.Bill `json:"bills"`
}
