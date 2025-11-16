package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/ran-jita/billing-engine/internal/model"
	"time"
)

type PaymentBillingRepository struct {
	db *sqlx.DB
}

func NewPaymentBillingRepository(db *sqlx.DB) *PaymentBillingRepository {
	return &PaymentBillingRepository{db: db}
}

// CreatePaymentBilling create new billing payment
func (r *PaymentBillingRepository) CreatePaymentBilling(ctx context.Context, billingPayment *model.PaymentBilling) error {
	query := `
       INSERT INTO payment_billings (
                             id,  
                             loan_id, 
                             billing_id,
                             amount,
                             created_at, 
                             updated_at
	 	) VALUES ($1, $2, $3, $4, $5, $6)
       RETURNING id, created_at, updated_at
   `
	billingPayment.ID = uuid.New().String()
	billingPayment.CreatedAt = time.Now()
	billingPayment.UpdatedAt = time.Now()

	err := r.db.QueryRowContext(
		ctx,
		query,
		billingPayment.ID,
		billingPayment.PaymentId,
		billingPayment.BillingId,
		billingPayment.Amount,
		billingPayment.CreatedAt,
		billingPayment.UpdatedAt,
	).Scan(&billingPayment.ID, &billingPayment.CreatedAt, &billingPayment.UpdatedAt)

	return err
}

// GetByID get loan by ID
//func (r *BorrowerRepository) ByID(ctx context.Context, id string) (model.Borrower, error) {
//	var borrower model.Borrower
//	query := `
//       SELECT
//           	id,
//			name,
//			delinquent,
//			created_at,
//			updated_at
//       FROM borrowers
//       WHERE id = $1
//   `
//
//	err := r.db.GetContext(ctx, &borrower, query, id)
//	if err != nil {
//		return borrower, err
//	}
//
//	return borrower, nil
//}

//// GetAll get all loan based on borrowerId
//func (r *LoanRepository) GetAll(ctx context.Context, borrowerId string) ([]model.Loan, error) {
//	var loans []model.Loan
//	query := `
//        SELECT
//            id,
//            borrower_id,
//            base_amount,
//            fee_amount,
//        	total_loan_amount,
//            outstanding_amount,
//            start_date,
//            created_at,
//            updated_at
//        FROM loans
//        WHERE borrower_id = $1
//        ORDER BY created_at DESC
//    `
//
//	err := r.db.SelectContext(ctx, &loans, query, borrowerId)
//	if err != nil {
//		return nil, err
//	}
//
//	return loans, nil
//}

// Update memperbarui loan
//func (r *LoanRepository) Update(ctx context.Context, loan *domain.Loan) error {
//	query := `
//        UPDATE loans
//        SET borrower_id = $1, principal_amount = $2, interest_rate = $3,
//            term = $4, start_date = $5, status = $6, updated_at = $7
//        WHERE id = $8
//    `
//
//	loan.UpdatedAt = time.Now()
//
//	_, err := r.db.ExecContext(
//		ctx,
//		query,
//		loan.BorrowerID,
//		loan.PrincipalAmount,
//		loan.InterestRate,
//		loan.Term,
//		loan.StartDate,
//		loan.Status,
//		loan.UpdatedAt,
//		loan.ID,
//	)
//
//	return err
//}

// Delete menghapus loan
//func (r *LoanRepository) Delete(ctx context.Context, id string) error {
//	query := `DELETE FROM loans WHERE id = $1`
//	_, err := r.db.ExecContext(ctx, query, id)
//	return err
//}
