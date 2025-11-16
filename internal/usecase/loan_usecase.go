package usecase

import (
	"context"
	"github.com/ran-jita/billing-engine/internal/domain"
	"github.com/ran-jita/billing-engine/internal/model"
)

type LoanUsecase struct {
	loanDomain *domain.LoanDomain
}

func NewLoanUsecase(loanDomain *domain.LoanDomain) *LoanUsecase {
	return &LoanUsecase{loanDomain: loanDomain}
}

func (h *LoanUsecase) GetAll(ctx context.Context, borrowerId string) ([]model.Loan, error) {
	return h.loanDomain.GetAll(ctx, borrowerId)
}

func (h *LoanUsecase) GetById(ctx context.Context, loanId string) (model.Loan, error) {
	return h.loanDomain.GetById(ctx, loanId)
}
