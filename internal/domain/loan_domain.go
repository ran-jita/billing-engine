package domain

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"github.com/ran-jita/billing-engine/internal/model/dto"
	"github.com/ran-jita/billing-engine/internal/model/postgresql"
	"github.com/ran-jita/billing-engine/internal/repository"
)

type LoanDomain struct {
	loanRepository *repository.LoanRepository
	billRepository *repository.BillRepository
}

func NewLoanDomain(
	loanRepository *repository.LoanRepository,
	billRepository *repository.BillRepository,
) *LoanDomain {
	return &LoanDomain{
		loanRepository: loanRepository,
		billRepository: billRepository,
	}
}

func (h *LoanDomain) GetAll(ctx context.Context, borrowerId string) ([]postgresql.Loan, error) {
	loan, err := h.loanRepository.GetAll(ctx, borrowerId)
	if len(loan) == 0 {
		// Handle not found case
		err = errors.New("loan not found")
		return nil, err
	}
	return loan, err
}

func (h *LoanDomain) GetById(ctx context.Context, loanId string) (postgresql.Loan, error) {
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

func (h *LoanDomain) GetOverdueBillByLoanId(ctx context.Context, loanId string) (dto.LoanWithBills, error) {
	var (
		response dto.LoanWithBills
		err      error
	)
	response.Loan, err = h.loanRepository.GetByID(ctx, loanId)
	if err != nil {
		// Handle UUID parsing error from database
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "22P02" { // invalid_text_representation (invalid UUID)
				err = errors.New("loan not found")
				return response, err
			}
		}

		if err == sql.ErrNoRows {
			err = errors.New("loan not found")
			return response, err
		}
	}

	response.Bills, err = h.billRepository.GetAllOverdueByLoanId(ctx, loanId)
	if err != nil {
		return response, err
	}

	return response, err
}

func (h *LoanDomain) PayBills(ctx context.Context, loanWithBills dto.LoanWithBills) error {
	var (
		err    error
		loanId string = loanWithBills.Loan.ID
	)

	for _, bill := range loanWithBills.Bills {
		err = h.billRepository.UpdateBillToPaid(ctx, bill.ID)
		if err != nil {
			return err
		}
	}

	totalPaid, err := h.billRepository.GetTotalPaid(ctx, loanId)
	if err != nil {
		return err
	}

	err = h.loanRepository.UpdateOutstandingAmount(ctx, loanId, *totalPaid)
	if err != nil {
		return err
	}

	return nil
}
