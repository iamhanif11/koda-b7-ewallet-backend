CREATE TABLE transfer_detail (
    id SERIAL PRIMARY KEY,
    transaction_id INTEGER REFERENCES transactions(id),
    receiver_id INTEGER REFERENCES users(id),
    notes VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW()
);