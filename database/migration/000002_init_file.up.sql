CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    account_id INT NOT NULL,
    operation_type_id INT NOT NULL,
    amount FLOAT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
);
CREATE TABLE operation_type (
    id SERIAL PRIMARY KEY,
    description VARCHAR(255) NOT NULL UNIQUE,
    coefficient INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP DEFAULT NULL
);
CREATE TABLE account (
    id SERIAL PRIMARY KEY,
    document_number VARCHAR(255) NOT NULL UNIQUE,    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
INSERT INTO operation_type (description, coefficient)
VALUES
    ('Normal Purchase', -1),
    ('Purchase with installments', -1),
    ('Withdrawal', -1),
    ('Credit Voucher', 1)
ON CONFLICT (description) DO NOTHING;
