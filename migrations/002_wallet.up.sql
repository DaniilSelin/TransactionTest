CREATE TABLE IF NOT EXISTS {{.Schema}}.wallets (
    address TEXT PRIMARY KEY,
    balance DECIMAL(18, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);