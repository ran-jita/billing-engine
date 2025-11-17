package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/ran-jita/billing-engine/internal/model/postgresql"
	"time"
)

type PaymentRepository struct {
	db *sqlx.DB
}

func NewPaymentRepository(db *sqlx.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

// CreatePayment create new payment
func (r *PaymentRepository) CreatePayment(ctx context.Context, tx *sqlx.Tx, payment *postgresql.Payment) error {
	query := `
       INSERT INTO payments (
                             id,  
                             borrower_id,
                             total_amount, 
                             payment_date, 
                             created_at, 
                             updated_at
	 	) VALUES ($1, $2, $3, $4, $5, $6)
       RETURNING id, created_at, updated_at
   `
	payment.ID = uuid.New().String()
	payment.CreatedAt = time.Now()
	payment.UpdatedAt = time.Now()

	err := tx.QueryRowContext(
		ctx,
		query,
		payment.ID,
		payment.BorrowerID,
		payment.TotalAmount,
		payment.PaymentDate,
		payment.CreatedAt,
		payment.UpdatedAt,
	).Scan(&payment.ID, &payment.CreatedAt, &payment.UpdatedAt)

	return err
}
