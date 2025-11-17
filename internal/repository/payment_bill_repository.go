package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/ran-jita/billing-engine/internal/model/postgresql"
	"time"
)

type PaymentBillRepository struct {
	db *sqlx.DB
}

func NewPaymentBillRepository(db *sqlx.DB) *PaymentBillRepository {
	return &PaymentBillRepository{db: db}
}

// CreatePaymentBill create new payment for a bill
func (r *PaymentBillRepository) CreatePaymentBill(ctx context.Context, tx *sqlx.Tx, paymentBill *postgresql.PaymentBill) error {
	query := `
       INSERT INTO payment_bills (
                             id,  
                        	 borrower_id,
                             payment_id, 
                             bill_id,
                             amount,
                             created_at, 
                             updated_at
	 	) VALUES ($1, $2, $3, $4, $5, $6, $7)
       RETURNING id, created_at, updated_at
   `
	paymentBill.ID = uuid.New().String()
	paymentBill.CreatedAt = time.Now()
	paymentBill.UpdatedAt = time.Now()

	err := tx.QueryRowContext(
		ctx,
		query,
		paymentBill.ID,
		paymentBill.BorrowerID,
		paymentBill.PaymentId,
		paymentBill.BillId,
		paymentBill.Amount,
		paymentBill.CreatedAt,
		paymentBill.UpdatedAt,
	).Scan(&paymentBill.ID, &paymentBill.CreatedAt, &paymentBill.UpdatedAt)

	return err
}
