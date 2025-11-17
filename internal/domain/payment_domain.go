package domain

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/ran-jita/billing-engine/internal/model/postgresql"
	"github.com/ran-jita/billing-engine/internal/repository"
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

func (d *PaymentDomain) CreatePayment(ctx context.Context, tx *sqlx.Tx, payment *postgresql.Payment, bills []postgresql.Bill) error {
	var err error

	err = d.paymentRepository.CreatePayment(ctx, tx, payment)
	if err != nil {
		fmt.Println("failed to create payment: ", err)
		return err
	}

	for _, bill := range bills {
		var paymentBill *postgresql.PaymentBill
		paymentBill = &postgresql.PaymentBill{
			BorrowerID: payment.BorrowerID,
			PaymentId:  payment.ID,
			BillId:     bill.ID,
			Amount:     bill.Amount,
		}
		err = d.paymentBilingRepository.CreatePaymentBill(ctx, tx, paymentBill)
		if err != nil {
			fmt.Println("failed to create payment bill: ", err)
			return err
		}
	}

	return nil
}
