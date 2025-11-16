package usecase

import (
	"context"
	"github.com/ran-jita/billing-engine/internal/domain"
	"github.com/ran-jita/billing-engine/internal/model"
	"time"
)

type PaymentUsecase struct {
	paymentDomain *domain.PaymentDomain
}

func NewPaymentUsecase(paymentDomain *domain.PaymentDomain) *PaymentUsecase {
	return &PaymentUsecase{paymentDomain: paymentDomain}
}

func (h *PaymentUsecase) Create(ctx context.Context, payment *model.Payment) error {
	payment.PaymentDate = time.Now()
	return h.paymentDomain.Create(ctx, payment)
}
