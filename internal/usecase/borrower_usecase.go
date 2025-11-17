package usecase

import (
	"context"
	"fmt"
	"github.com/ran-jita/billing-engine/internal/domain"
	"github.com/ran-jita/billing-engine/internal/model/postgresql"
	"time"
)

type BorrowerUsecase struct {
	borrowerDomain *domain.BorrowerDomain
	loanDomain     *domain.LoanDomain
}

func NewBorrowerUsecase(
	borrowerDomain *domain.BorrowerDomain,
	loanDomain *domain.LoanDomain,
) *BorrowerUsecase {
	return &BorrowerUsecase{
		borrowerDomain: borrowerDomain,
		loanDomain:     loanDomain,
	}
}

func (h *BorrowerUsecase) GetById(ctx context.Context, borrowerId string) (postgresql.Borrower, error) {
	return h.borrowerDomain.GetById(ctx, borrowerId)
}

func (h *BorrowerUsecase) UpdateDelinquent(ctx context.Context, processDate time.Time) error {
	processDate = processDate.AddDate(0, 0, -14)
	fmt.Println(processDate)
	borrowerIds, err := h.loanDomain.GetBorrowerIdWithOverdueBill(ctx, processDate)
	if err != nil {
		fmt.Println("borrowerIds error", err)
		return err
	}

	for _, borrowerId := range borrowerIds {
		err = h.borrowerDomain.UpdateDelinquentByBorrowerId(ctx, borrowerId)
		if err != nil {
			fmt.Println("update borrower error", err)
			return err
		}
	}

	return nil

}
