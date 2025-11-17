package domain

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/ran-jita/billing-engine/internal/model/dto"
	"github.com/ran-jita/billing-engine/internal/model/postgresql"
	"github.com/ran-jita/billing-engine/internal/repository"
	"time"
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

func (d *LoanDomain) GetAll(ctx context.Context, borrowerId string) ([]postgresql.Loan, error) {
	loan, err := d.loanRepository.GetAllByBorrowerId(ctx, borrowerId)
	if len(loan) == 0 {
		// Handle not found case
		err = errors.New("loan not found")
		return nil, err
	}
	return loan, err
}

func (d *LoanDomain) GetById(ctx context.Context, loanId string) (postgresql.Loan, error) {
	loan, err := d.loanRepository.GetByID(ctx, loanId)
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

func (d *LoanDomain) GetOverdueBillByLoanId(ctx context.Context, loanId string, paymentDate time.Time) (dto.LoanWithBills, error) {
	var (
		response dto.LoanWithBills
		err      error
	)
	response.Loan, err = d.loanRepository.GetByID(ctx, loanId)
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

	response.Bills, err = d.billRepository.GetAllOverdueByLoanId(ctx, loanId, paymentDate)
	if err != nil {
		return response, err
	}

	return response, err
}

func (d *LoanDomain) PayBills(ctx context.Context, tx *sqlx.Tx, loanWithBills dto.LoanWithBills) error {
	var (
		err    error
		loanId string = loanWithBills.Loan.ID
	)

	for _, bill := range loanWithBills.Bills {
		err = d.billRepository.UpdateBillToPaid(ctx, tx, bill.ID)
		if err != nil {
			return err
		}
	}

	totalPaid, err := d.billRepository.GetTotalPaid(ctx, tx, loanId)
	if err != nil {
		return err
	}

	err = d.loanRepository.UpdateOutstandingAmount(ctx, tx, loanId, *totalPaid)
	if err != nil {
		return err
	}

	return nil
}

func (d *LoanDomain) GetBorrowerIdWithOverdueBill(ctx context.Context, progressDate time.Time) ([]string, error) {
	var (
		borrowerId []string
		err        error
	)

	borrowerId, err = d.billRepository.GetBorrowerIdWithOverdueBills(ctx, progressDate)
	if err != nil {
		return borrowerId, err
	}

	return borrowerId, nil
}
