package repository

import (
	"context"
	"github.com/ran-jita/billing-engine/internal/model/postgresql"
	"time"

	//"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type LoanRepository struct {
	db *sqlx.DB
}

func NewLoanRepository(db *sqlx.DB) *LoanRepository {
	return &LoanRepository{db: db}
}

// GetByID get loan by ID
func (r *LoanRepository) GetByID(ctx context.Context, id string) (postgresql.Loan, error) {
	var loan postgresql.Loan
	query := `
       SELECT 
           id, 
            borrower_id, 
            base_amount, 
            fee_amount, 
        	total_loan_amount, 
            outstanding_amount,
            start_date, 
            created_at, 
            updated_at
       FROM loans
       WHERE id = $1
   `

	err := r.db.GetContext(ctx, &loan, query, id)
	if err != nil {
		return loan, err
	}

	return loan, nil
}

// GetAllByBorrowerId get all loan based on borrowerId
func (r *LoanRepository) GetAllByBorrowerId(ctx context.Context, borrowerId string) ([]postgresql.Loan, error) {
	var loans []postgresql.Loan
	query := `
        SELECT 
            id, 
            borrower_id, 
            base_amount, 
            fee_amount, 
        	total_loan_amount, 
            outstanding_amount,
            start_date, 
            created_at, 
            updated_at
        FROM loans
        WHERE borrower_id = $1
        ORDER BY created_at DESC
    `

	err := r.db.SelectContext(ctx, &loans, query, borrowerId)
	if err != nil {
		return nil, err
	}

	return loans, nil
}

// UpdateOutstandingAmount update outstanding amount based on total paid bill
func (r *LoanRepository) UpdateOutstandingAmount(ctx context.Context, tx *sqlx.Tx, loanId string, totalPaid float64) error {
	query := `
       UPDATE loans
       SET outstanding_amount = total_loan_amount - $1, updated_at = $2
       WHERE id = $3
   `

	updatedAt := time.Now()

	_, err := tx.ExecContext(
		ctx,
		query,
		totalPaid,
		updatedAt,
		loanId,
	)

	return err
}
