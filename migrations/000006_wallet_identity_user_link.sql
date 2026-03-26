-- +goose Up
ALTER TABLE auth_wallet_identities
ADD COLUMN IF NOT EXISTS user_id TEXT NULL REFERENCES users(id);

CREATE UNIQUE INDEX IF NOT EXISTS idx_auth_wallet_identities_user_id
ON auth_wallet_identities(user_id)
WHERE user_id IS NOT NULL;

-- +goose Down
DROP INDEX IF EXISTS idx_auth_wallet_identities_user_id;

ALTER TABLE auth_wallet_identities
DROP COLUMN IF EXISTS user_id;
