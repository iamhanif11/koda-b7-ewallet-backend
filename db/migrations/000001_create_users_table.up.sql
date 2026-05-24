CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    fullname VARCHAR(255),
    password VARCHAR(255) NOT NULL,
    pin CHAR(6),
    picture VARCHAR(255),
    phone VARCHAR(20),
    created_at TIMESTAMP DEFAULt NOW(),
    updated_at TIMESTAMP 
);