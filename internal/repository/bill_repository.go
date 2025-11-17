package repository

import (
	//"github.com/google/uuid"
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/ran-jita/billing-engine/internal/model/postgresql"
	"time"
)

type BillRepository struct {
	db *sqlx.DB
}

const statusBillUnpaid string = "unpaid"
const statusBillPaid string = "paid"

func NewBillRepository(db *sqlx.DB) *BillRepository {
	return &BillRepository{db: db}
}

// GetTotalPaid get sum of amount of bill with status paid
func (r *BillRepository) GetTotalPaid(ctx context.Context, tx *sqlx.Tx, loanId string) (*float64, error) {
	var totalPaid *float64
	query := `
      SELECT sum(amount)
		FROM bills
		WHERE loan_id = $1 AND status = $2
	`

	err := tx.GetContext(ctx, &totalPaid, query, loanId, statusBillPaid)
	if err != nil {
		return nil, err
	}

	return totalPaid, nil
}

// GetAllOverdueByLoanId get all overdue bills by loan_id
func (r *BillRepository) GetAllOverdueByLoanId(ctx context.Context, loanId string, paymentDate time.Time) ([]postgresql.Bill, error) {
	var bills []postgresql.Bill
	query := `
       SELECT
           id,
           borrower_id,
           loan_id,
           amount,
           due_date,
       	   status,
           created_at,
           updated_at
       FROM bills
       WHERE loan_id = $1 AND due_date <= $2 AND status = $3
       ORDER BY created_at DESC
   `

	err := r.db.SelectContext(ctx, &bills, query, loanId, paymentDate, statusBillUnpaid)
	if err != nil {
		return nil, err
	}

	return bills, nil
}

// GetBorrowerIdWithOverdueBills get borrower_id with overdue bills
func (r *BillRepository) GetBorrowerIdWithOverdueBills(ctx context.Context, processDate time.Time) ([]string, error) {
	var borrowerId []string
	query := `
       	SELECT borrower_id
		FROM bills
		WHERE due_date < $1 AND status = $2
		GROUP BY borrower_id
    	HAVING count(id) >= 2
   `

	err := r.db.SelectContext(ctx, &borrowerId, query, processDate, statusBillUnpaid)
	if err != nil {
		return nil, err
	}

	return borrowerId, nil
}

// UpdateBillToPaid update bills data
func (r *BillRepository) UpdateBillToPaid(ctx context.Context, tx *sqlx.Tx, billId string) error {
	query := `
       UPDATE bills
       SET status = $1, updated_at = $2
       WHERE id = $3
   `

	updatedAt := time.Now()
	_, err := tx.ExecContext(
		ctx,
		query,
		statusBillPaid,
		updatedAt,
		billId,
	)

	return err
}
