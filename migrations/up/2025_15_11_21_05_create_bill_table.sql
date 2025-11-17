CREATE TABLE bills (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    borrower_id VARCHAR(36) NOT NULL,
    loan_id VARCHAR(36) NOT NULL,
    amount FLOAT NOT NULL,
    due_date DATE NOT NULL,
    status VARCHAR(10) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);