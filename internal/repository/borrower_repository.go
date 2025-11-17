package repository

import (
	"context"
	"github.com/ran-jita/billing-engine/internal/model/postgresql"
	"time"

	//"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type BorrowerRepository struct {
	db *sqlx.DB
}

func NewBorrowerRepository(db *sqlx.DB) *BorrowerRepository {
	return &BorrowerRepository{db: db}
}

const statusUserIsDelinquent = true

// GetByID get borrower by ID
func (r *BorrowerRepository) GetByID(ctx context.Context, id string) (postgresql.Borrower, error) {
	var borrower postgresql.Borrower
	query := `
       SELECT 
           	id, 
			name, 
			delinquent, 
			created_at, 
			updated_at
       FROM borrowers
       WHERE id = $1
   `

	err := r.db.GetContext(ctx, &borrower, query, id)
	if err != nil {
		return borrower, err
	}

	return borrower, nil
}

// UpdateDelinquentStatus update borrower delinquent status
func (r *BorrowerRepository) UpdateDelinquentStatus(ctx context.Context, borrowerId string) error {
	query := `
       UPDATE borrowers
       SET delinquent = $1, updated_at = $2
       WHERE id = $3
   `

	updatedAt := time.Now()

	_, err := r.db.ExecContext(
		ctx,
		query,
		statusUserIsDelinquent,
		updatedAt,
		borrowerId,
	)

	return err
}
