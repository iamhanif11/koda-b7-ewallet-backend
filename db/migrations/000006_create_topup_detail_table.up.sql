CREATE TABLE topup_detail (
    id SERIAL PRIMARY KEY,
    transaction_id INTEGER REFERENCES transactions(id),
    payment_method_id INTEGER REFERENCES payment_method(id),
    service_fee INTEGER,
    tax_amount INTEGER,
    sub_total INTEGER
);