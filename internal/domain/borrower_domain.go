package domain

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"github.com/ran-jita/billing-engine/internal/model"
	"github.com/ran-jita/billing-engine/internal/repository"
)

type BorrowerDomain struct {
	borrowerRepository *repository.BorrowerRepository
}

func NewBorrowerDomain(borrowerRepository *repository.BorrowerRepository) *BorrowerDomain {
	return &BorrowerDomain{borrowerRepository: borrowerRepository}
}

func (h *BorrowerDomain) GetById(ctx context.Context, borrowerId string) (model.Borrower, error) {
	borrower, err := h.borrowerRepository.GetByID(ctx, borrowerId)
	if err != nil {
		// Handle UUID parsing error from database
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "22P02" { // invalid_text_representation (invalid UUID)
				err = errors.New("loan not found")
				return borrower, err
			}
		}

		if err == sql.ErrNoRows {
			err = errors.New("loan not found")
			return borrower, err
		}
	}
	return borrower, err
}
