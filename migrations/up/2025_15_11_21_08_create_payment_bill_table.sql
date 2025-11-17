CREATE TABLE payment_bills (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    borrower_id VARCHAR(36) NOT NULL,
    payment_id VARCHAR(36) NOT NULL,
    bill_id VARCHAR(36) NOT NULL,
    amount FLOAT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);