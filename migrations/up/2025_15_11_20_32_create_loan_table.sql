CREATE TABLE loans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    borrower_id VARCHAR(36) NOT NULL,
    base_amount FLOAT NOT NULL,
    fee_amount FLOAT NOT NULL,
    total_loan_amount FLOAT NOT NULL,
    outstanding_amount FLOAT NOT NULL,
    start_date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);