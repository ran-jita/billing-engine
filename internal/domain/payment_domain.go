package domain

import (
	"context"
	"github.com/ran-jita/billing-engine/internal/model/postgresql"
	"github.com/ran-jita/billing-engine/internal/repository"
	"time"
)

type PaymentDomain struct {
	paymentRepository       *repository.PaymentRepository
	paymentBilingRepository *repository.PaymentBillRepository
}

func NewPaymentDomain(
	paymentRepository *repository.PaymentRepository,
	paymentBillRepository *repository.PaymentBillRepository,
) *PaymentDomain {
	return &PaymentDomain{
		paymentRepository:       paymentRepository,
		paymentBilingRepository: paymentBillRepository,
	}
}

func (h *PaymentDomain) CreatePayment(ctx context.Context, payment *postgresql.Payment, bills []postgresql.Bill) error {
	var err error

	payment.PaymentDate = time.Now()
	err = h.paymentRepository.CreatePayment(ctx, payment)
	if err != nil {
		return err
	}

	for _, bill := range bills {
		var paymentBill *postgresql.PaymentBill
		paymentBill = &postgresql.PaymentBill{
			PaymentId: payment.ID,
			BillId:    bill.ID,
			Amount:    bill.Amount,
		}
		err = h.paymentBilingRepository.CreatePaymentBill(ctx, paymentBill)
		if err != nil {
			return err
		}
	}

	return nil
}
