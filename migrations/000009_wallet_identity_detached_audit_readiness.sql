-- +goose Up
ALTER TABLE auth_wallet_identities
ADD COLUMN IF NOT EXISTS detached_at TIMESTAMPTZ NULL;

CREATE INDEX IF NOT EXISTS idx_auth_wallet_identities_detached_at
ON auth_wallet_identities(detached_at)
WHERE detached_at IS NOT NULL;

-- +goose Down
DROP INDEX IF EXISTS idx_auth_wallet_identities_detached_at;

ALTER TABLE auth_wallet_identities
DROP COLUMN IF EXISTS detached_at;
