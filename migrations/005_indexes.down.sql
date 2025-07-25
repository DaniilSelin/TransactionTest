DROP INDEX IF EXISTS {{.Schema}}.idx_transactions_to_wallet_created_at;

DROP INDEX IF EXISTS {{.Schema}}.idx_transactions_from_wallet_created_at;

DROP INDEX IF EXISTS {{.Schema}}.idx_transactions_to_wallet;

DROP INDEX IF EXISTS {{.Schema}}.idx_transactions_from_wallet;
