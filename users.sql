CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    email VARCHAR(100) UNIQUE NOT NULL,
    balance NUMERIC(10,2) DEFAULT 0
);

INSERT INTO users (name, email, balance) VALUES
('Gojo', 'gojo@example.com', 100.00),
('Eren', 'eren@example.com', 50.00),
('Naruto', 'naruto@example.com', 75.00);

