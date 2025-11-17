package usecase

import (
	"context"
	"github.com/ran-jita/billing-engine/internal/domain"
	"github.com/ran-jita/billing-engine/internal/model/postgresql"
)

type LoanUsecase struct {
	loanDomain *domain.LoanDomain
}

func NewLoanUsecase(loanDomain *domain.LoanDomain) *LoanUsecase {
	return &LoanUsecase{loanDomain: loanDomain}
}

func (u *LoanUsecase) GetAll(ctx context.Context, borrowerId string) ([]postgresql.Loan, error) {
	return u.loanDomain.GetAll(ctx, borrowerId)
}

func (u *LoanUsecase) GetById(ctx context.Context, loanId string) (postgresql.Loan, error) {
	return u.loanDomain.GetById(ctx, loanId)
}
