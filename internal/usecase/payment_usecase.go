package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/ran-jita/billing-engine/internal/domain"
	"github.com/ran-jita/billing-engine/internal/model/dto"
	"github.com/ran-jita/billing-engine/internal/model/postgresql"
	"time"
)

type PaymentUsecase struct {
	paymentDomain *domain.PaymentDomain
	loanDomain    *domain.LoanDomain
	db            *sqlx.DB
}

func NewPaymentUsecase(
	paymentDomain *domain.PaymentDomain,
	loanDomain *domain.LoanDomain,
	db *sqlx.DB,
) *PaymentUsecase {
	return &PaymentUsecase{
		paymentDomain: paymentDomain,
		loanDomain:    loanDomain,
		db:            db,
	}
}

func (u *PaymentUsecase) Create(ctx context.Context, request *dto.CreatePayment) (postgresql.Payment, error) {
	var (
		payment        postgresql.Payment
		loanWithBills  dto.LoanWithBills
		paymentDate    time.Time
		totalAmountDue float64
		totalPayment   float64
		err            error
	)

	paymentDate = time.Now()
	if request.PaymentDate != "" {
		paymentDate, _ = time.Parse("2006-01-02", request.PaymentDate)
	}
	totalPayment = request.TotalAmount

	loanWithBills, err = u.loanDomain.GetOverdueBillByLoanId(ctx, request.LoanId, paymentDate)
	if err != nil {
		fmt.Println("error getting loan with bills: ", err)
		return payment, err
	}

	if len(loanWithBills.Bills) == 0 {
		err = errors.New("no overdue bills found")
		return payment, err
	}

	for _, bill := range loanWithBills.Bills {
		totalAmountDue += bill.Amount
	}

	if totalAmountDue != request.TotalAmount {
		err = errors.New("total amount is not match with overdue bill")
		return payment, err
	}

	tx, err := u.db.Beginx()
	if err != nil {
		fmt.Println("error begin transaction: ", err)
		return payment, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Prepare payment data
	payment.BorrowerID = loanWithBills.Loan.BorrowerID
	payment.PaymentDate = paymentDate
	payment.TotalAmount = totalPayment

	err = u.paymentDomain.CreatePayment(ctx, tx, &payment, loanWithBills.Bills)
	if err != nil {
		fmt.Println("error create payment: ", err)
		return payment, fmt.Errorf("create payment failed: %w", err)
	}

	err = u.loanDomain.PayBills(ctx, tx, loanWithBills)
	if err != nil {
		fmt.Println("error paying bills: ", err)
		return payment, fmt.Errorf("update bill failed: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		fmt.Println("error commit transaction: ", err)
		return payment, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return payment, nil
}
