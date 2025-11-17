package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/ran-jita/billing-engine/internal/domain"
	"github.com/ran-jita/billing-engine/internal/model/dto"
	"github.com/ran-jita/billing-engine/internal/model/postgresql"
	"sync"
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
		return payment, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var wg sync.WaitGroup
	errChan := make(chan error, 2)

	wg.Add(2)

	go func() {
		defer wg.Done()
		payment.PaymentDate = paymentDate
		payment.TotalAmount = totalPayment
		err = u.paymentDomain.CreatePayment(ctx, tx, &payment, loanWithBills.Bills)
		if err != nil {
			errChan <- fmt.Errorf("create payment failed: %w", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := u.loanDomain.PayBills(ctx, tx, loanWithBills); err != nil {
			errChan <- fmt.Errorf("update bill failed: %w", err)
		}
	}()

	wg.Wait()
	close(errChan)

	// Collect all errors from channel
	var errs []error
	for err := range errChan {
		errs = append(errs, err)
	}

	// If there are any errors, rollback and return
	if len(errs) > 0 {
		return payment, fmt.Errorf("payment processing failed: %v", errs)
	}

	err = tx.Commit()
	if err != nil {
		return payment, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return payment, nil
}
