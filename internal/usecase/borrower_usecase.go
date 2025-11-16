package usecase

import (
	"context"
	"github.com/ran-jita/billing-engine/internal/domain"
	"github.com/ran-jita/billing-engine/internal/model"
)

type BorrowerUsecase struct {
	borrowerDomain *domain.BorrowerDomain
}

func NewBorrowerUsecase(borrowerDomain *domain.BorrowerDomain) *BorrowerUsecase {
	return &BorrowerUsecase{borrowerDomain: borrowerDomain}
}

func (h *BorrowerUsecase) GetById(ctx context.Context, borrowerId string) (model.Borrower, error) {
	return h.borrowerDomain.GetById(ctx, borrowerId)
}
