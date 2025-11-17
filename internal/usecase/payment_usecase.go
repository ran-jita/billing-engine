package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/ran-jita/billing-engine/internal/domain"
	"github.com/ran-jita/billing-engine/internal/model/dto"
	"github.com/ran-jita/billing-engine/internal/model/postgresql"
	"sync"
)

type PaymentUsecase struct {
	paymentDomain *domain.PaymentDomain
	loanDomain    *domain.LoanDomain
}

func NewPaymentUsecase(
	paymentDomain *domain.PaymentDomain,
	loanDomain *domain.LoanDomain,
) *PaymentUsecase {
	return &PaymentUsecase{
		paymentDomain: paymentDomain,
		loanDomain:    loanDomain,
	}
}

func (u *PaymentUsecase) Create(ctx context.Context, request *dto.CreatePayment) (postgresql.Payment, error) {
	var (
		payment        postgresql.Payment
		loanWithBills  dto.LoanWithBills
		totalAmountDue float64
		err            error
	)

	loanWithBills, err = u.loanDomain.GetOverdueBillByLoanId(ctx, request.LoandId)
	if err != nil {
		return payment, err
	}

	for _, bill := range loanWithBills.Bills {
		totalAmountDue += bill.Amount
	}

	if totalAmountDue != request.TotalAmount {
		err = errors.New("total amount is wrong")
		return payment, err
	}

	var wg sync.WaitGroup
	errChan := make(chan error, 2) // Buffer 2 untuk 2 goroutine

	wg.Add(2)

	go func() {
		defer wg.Done()
		err = u.paymentDomain.CreatePayment(ctx, &payment, loanWithBills.Bills)
		if err != nil {
			errChan <- fmt.Errorf("create payment failed: %w", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := u.loanDomain.PayBills(ctx, loanWithBills); err != nil {
			errChan <- fmt.Errorf("update bill failed: %w", err)
		}
	}()

	wg.Wait()
	close(errChan)

	var errs []error
	for err := range errChan {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return payment, fmt.Errorf("payment processing failed: %v", errs)
	}

	return payment, <-errChan
}
