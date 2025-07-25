ALTER TABLE {{.Schema}}.wallets
  ADD CONSTRAINT chk_balance_nonnegative CHECK (balance >= 0);

ALTER TABLE {{.Schema}}.transactions
  ADD CONSTRAINT chk_amount_positive CHECK (amount > 0),
  ADD CONSTRAINT chk_no_self_transfer CHECK (from_wallet <> to_wallet);