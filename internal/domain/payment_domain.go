package domain

import (
	"context"
	"github.com/ran-jita/billing-engine/internal/model"
	"github.com/ran-jita/billing-engine/internal/repository"
)

type PaymentDomain struct {
	paymentRepository *repository.PaymentRepository
}

func NewPaymentDomain(paymentRepository *repository.PaymentRepository) *PaymentDomain {
	return &PaymentDomain{paymentRepository: paymentRepository}
}

func (h *PaymentDomain) Create(ctx context.Context, payment *model.Payment) error {
	err := h.paymentRepository.Create(ctx, payment)
	if err != nil {
		return err
	}
	return nil
}
