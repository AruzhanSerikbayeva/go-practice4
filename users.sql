CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
   name TEXT NOT NULL,
   email TEXT UNIQUE NOT NULL,
   balance NUMERIC DEFAULT 0
);
-- просто для теста
INSERT INTO users (name, email, balance) VALUES
('Gojo', 'gojo@example.com', 100.00),
('Eren', 'eren@example.com', 50.00),
('Naruto', 'naruto@example.com', 75.00);

