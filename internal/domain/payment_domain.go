package domain

import (
	"context"
	"github.com/ran-jita/billing-engine/internal/model"
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

func (h *PaymentDomain) CreatePayment(ctx context.Context, payment *model.Payment, bills []model.Bill) error {
	var err error

	payment.PaymentDate = time.Now()
	err = h.paymentRepository.CreatePayment(ctx, payment)
	if err != nil {
		return err
	}

	for _, bill := range bills {
		var paymentBill *model.PaymentBill
		paymentBill = &model.PaymentBill{
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
