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

// Create membuat loan baru
//func (r *LoanRepository) Create(ctx context.Context, loan *domain.Loan) error {
//	query := `
//        INSERT INTO loans (id, borrower_id, principal_amount, interest_rate, term, start_date, status, created_at, updated_at)
//        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
//        RETURNING id, created_at, updated_at
//    `
//
//	loan.ID = uuid.New().String()
//	loan.CreatedAt = time.Now()
//	loan.UpdatedAt = time.Now()
//
//	err := r.db.QueryRowContext(
//		ctx,
//		query,
//		loan.ID,
//		loan.BorrowerID,
//		loan.BaseAmount,
//		loan.InterestRate,
//		loan.Term,
//		loan.StartDate,
//		loan.Status,
//		loan.CreatedAt,
//		loan.UpdatedAt,
//	).Scan(&loan.ID, &loan.CreatedAt, &loan.UpdatedAt)
//
//	return err
//}

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

// GetAll get all loan based on borrowerId
func (r *LoanRepository) GetAll(ctx context.Context, borrowerId string) ([]postgresql.Loan, error) {
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

// Update update outstanding amount based on total paid bill
func (r *LoanRepository) UpdateOutstandingAmount(ctx context.Context, tx *sqlx.Tx, loanId string, totalPaid float64) error {
	query := `
       UPDATE loans
       SET outstanding_amount = total_loan_amount - $1, updated_at = $2
       WHERE id = $3
   `

	_, err := tx.ExecContext(
		ctx,
		query,
		totalPaid,
		time.Now(),
		loanId,
	)

	return err
}

// Delete menghapus loan
//func (r *LoanRepository) Delete(ctx context.Context, id string) error {
//	query := `DELETE FROM loans WHERE id = $1`
//	_, err := r.db.ExecContext(ctx, query, id)
//	return err
//}
