package dto

import (
	"github.com/ran-jita/billing-engine/internal/model/postgresql"
)

type LoanWithBills struct {
	Loan  postgresql.Loan   `json:"loan"`
	Bills []postgresql.Bill `json:"bills"`
}
