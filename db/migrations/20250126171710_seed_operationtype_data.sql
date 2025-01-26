-- +goose Up
INSERT INTO operation_type (description, coefficient)
VALUES
    ('Normal Purchase', -1),
    ('Purchase with installments', -1),
    ('Withdrawal', -1),
    ('Credit Voucher', 1)
ON CONFLICT (description) DO NOTHING;

-- +goose Down
DELETE FROM operation_type
WHERE description IN ('Normal Purchase', 'Purchase with installments', 'Withdrawal', 'Credit Voucher');
