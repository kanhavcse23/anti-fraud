DROP TABLE transactions;
DROP TABLE operation_types;
DROP TABLE account;
DELETE FROM operation_type
    WHERE description IN ('Normal Purchase', 'Purchase with installments', 'Withdrawal', 'Credit Voucher');