CREATE TABLE IF NOT EXISTS {{.Schema}}.transactions (
    id SERIAL PRIMARY KEY,
    from_wallet TEXT NOT NULL,
    to_wallet TEXT NOT NULL,
    amount DECIMAL(18, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    CONSTRAINT fk_from FOREIGN KEY (from_wallet) REFERENCES {{.Schema}}.wallets(address),
    CONSTRAINT fk_to FOREIGN KEY (to_wallet) REFERENCES {{.Schema}}.wallets(address)
);

CREATE INDEX IF NOT EXISTS idx_created_at ON {{.Schema}}.transactions (created_at DESC);