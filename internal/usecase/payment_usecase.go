package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/ran-jita/billing-engine/internal/domain"
	"github.com/ran-jita/billing-engine/internal/model"
	"github.com/ran-jita/billing-engine/internal/model/dto"
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

func (u *PaymentUsecase) Create(ctx context.Context, request *dto.CreatePayment) (model.Payment, error) {
	var (
		payment          model.Payment
		loanWithBillings dto.LoanWithBillings
		totalAmountDue   float64
		err              error
	)

	loanWithBillings, err = u.loanDomain.GetOverdueBillingByLoanId(ctx, request.LoandId)
	if err != nil {
		return payment, err
	}

	for _, billing := range loanWithBillings.Billings {
		totalAmountDue += billing.Amount
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
		err = u.paymentDomain.CreatePayment(ctx, &payment, loanWithBillings.Billings)
		if err != nil {
			errChan <- fmt.Errorf("create payment failed: %w", err)
		}
	}()

	// Goroutine 2: Update billing
	go func() {
		defer wg.Done()
		if err := u.loanDomain.PayBillings(ctx, loanWithBillings); err != nil {
			errChan <- fmt.Errorf("update billing failed: %w", err)
		}
	}()

	wg.Wait()
	close(errChan)

	var errs []error
	for err := range errChan {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return fmt.Errorf("payment processing failed: %v", errs)
	}

	return payment, <-errChan
}
