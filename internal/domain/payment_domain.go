package domain

import (
	"context"
	"github.com/ran-jita/billing-engine/internal/model"
	"github.com/ran-jita/billing-engine/internal/repository"
	"time"
)

type PaymentDomain struct {
	paymentRepository       *repository.PaymentRepository
	paymentBilingRepository *repository.PaymentBillingRepository
}

func NewPaymentDomain(
	paymentRepository *repository.PaymentRepository,
	paymentBillingRepository *repository.PaymentBillingRepository,
) *PaymentDomain {
	return &PaymentDomain{
		paymentRepository:       paymentRepository,
		paymentBilingRepository: paymentBillingRepository,
	}
}

func (h *PaymentDomain) CreatePayment(ctx context.Context, payment *model.Payment, billings []model.Billing) error {
	var err error

	payment.PaymentDate = time.Now()
	err = h.paymentRepository.CreatePayment(ctx, payment)
	if err != nil {
		return err
	}

	for _, billing := range billings {
		var paymentBilling *model.PaymentBilling
		paymentBilling = &model.PaymentBilling{
			PaymentId: payment.ID,
			BillingId: billing.ID,
			Amount:    billing.Amount,
		}
		err = h.paymentBilingRepository.CreatePaymentBilling(ctx, paymentBilling)
		if err != nil {
			return err
		}
	}

	return nil
}
