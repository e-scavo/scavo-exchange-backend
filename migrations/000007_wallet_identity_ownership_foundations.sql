-- +goose Up
DROP INDEX IF EXISTS idx_auth_wallet_identities_user_id;

ALTER TABLE auth_wallet_identities
ADD COLUMN IF NOT EXISTS linked_at TIMESTAMPTZ NULL,
ADD COLUMN IF NOT EXISTS is_primary BOOLEAN NOT NULL DEFAULT FALSE;

UPDATE auth_wallet_identities
SET
    linked_at = COALESCE(linked_at, NOW()),
    is_primary = TRUE
WHERE user_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_auth_wallet_identities_user_id
ON auth_wallet_identities(user_id)
WHERE user_id IS NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_auth_wallet_identities_user_primary
ON auth_wallet_identities(user_id)
WHERE user_id IS NOT NULL AND is_primary = TRUE;

-- +goose Down
DROP INDEX IF EXISTS idx_auth_wallet_identities_user_primary;
DROP INDEX IF EXISTS idx_auth_wallet_identities_user_id;

ALTER TABLE auth_wallet_identities
DROP COLUMN IF EXISTS is_primary,
DROP COLUMN IF EXISTS linked_at;

CREATE UNIQUE INDEX IF NOT EXISTS idx_auth_wallet_identities_user_id
ON auth_wallet_identities(user_id)
WHERE user_id IS NOT NULL;