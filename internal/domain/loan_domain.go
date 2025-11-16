package domain

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"github.com/ran-jita/billing-engine/internal/model"
	"github.com/ran-jita/billing-engine/internal/repository"
)

type LoanDomain struct {
	loanRepository *repository.LoanRepository
}

func NewLoanDomain(loanRepository *repository.LoanRepository) *LoanDomain {
	return &LoanDomain{loanRepository: loanRepository}
}

func (h *LoanDomain) GetAll(ctx context.Context, borrowerId string) ([]model.Loan, error) {
	loan, err := h.loanRepository.GetAll(ctx, borrowerId)
	if len(loan) == 0 {
		// Handle not found case
		err = errors.New("loan not found")
		return nil, err
	}
	return loan, err
}

func (h *LoanDomain) GetById(ctx context.Context, loanId string) (model.Loan, error) {
	loan, err := h.loanRepository.GetByID(ctx, loanId)
	if err != nil {
		// Handle UUID parsing error from database
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "22P02" { // invalid_text_representation (invalid UUID)
				err = errors.New("loan not found")
				return loan, err
			}
		}

		if err == sql.ErrNoRows {
			err = errors.New("loan not found")
			return loan, err
		}
	}
	return loan, err
}
