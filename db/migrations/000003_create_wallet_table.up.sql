CREATE TABLE wallet (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    balance INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW (),
    update_at TIMESTAMP 
);