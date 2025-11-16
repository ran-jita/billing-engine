package repository

import (
	//"github.com/google/uuid"
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/ran-jita/billing-engine/internal/model"
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

// GetPaid get sum of amount of bill with status paid
func (r *BillRepository) GetTotalPaid(ctx context.Context, loanId string) (*float64, error) {
	var totalPaid *float64
	query := `
      SELECT sum(amount)
		FROM bills
		WHERE loan_id = $1 AND status = $2
	`

	err := r.db.GetContext(ctx, &totalPaid, query, loanId, statusBillPaid)
	if err != nil {
		return nil, err
	}

	return totalPaid, nil
}

// GetAllOverdueByLoanId get all overdue bills by loan_id
func (r *BillRepository) GetAllOverdueByLoanId(ctx context.Context, loanId string) ([]model.Bill, error) {
	var bills []model.Bill
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
       WHERE id = $1 AND due_date <= NOW() AND status = $2
       ORDER BY created_at DESC
   `

	err := r.db.SelectContext(ctx, &bills, query, loanId, statusBillUnpaid)
	if err != nil {
		return nil, err
	}

	return bills, nil
}

// UpdateBill update bills data
func (r *BillRepository) UpdateBillToPaid(ctx context.Context, billId string) error {
	query := `
       UPDATE bills
       SET status = $1, updated_at = $2
       WHERE id = $3
   `

	_, err := r.db.ExecContext(
		ctx,
		query,
		statusBillPaid,
		time.Now,
		billId,
	)

	return err
}

// Delete menghapus loan
//func (r *LoanRepository) Delete(ctx context.Context, id string) error {
//	query := `DELETE FROM loans WHERE id = $1`
//	_, err := r.db.ExecContext(ctx, query, id)
//	return err
//}
