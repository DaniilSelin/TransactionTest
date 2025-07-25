ALTER TABLE {{.Schema}}.transactions
  DROP CONSTRAINT IF EXISTS chk_amount_positive,
  DROP CONSTRAINT IF EXISTS chk_no_self_transfer;

ALTER TABLE {{.Schema}}.wallets
  DROP CONSTRAINT IF EXISTS chk_balance_nonnegative;