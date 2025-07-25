CREATE INDEX IF NOT EXISTS idx_transactions_from_wallet
ON {{.Schema}}.transactions (from_wallet);

CREATE INDEX IF NOT EXISTS idx_transactions_to_wallet
ON {{.Schema}}.transactions (to_wallet);

CREATE INDEX IF NOT EXISTS idx_transactions_from_wallet_created_at
ON {{.Schema}}.transactions (from_wallet, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_transactions_to_wallet_created_at
ON {{.Schema}}.transactions (to_wallet, created_at DESC);
