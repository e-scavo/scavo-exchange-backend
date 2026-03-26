ALTER TABLE auth_wallet_challenges
ADD COLUMN IF NOT EXISTS purpose TEXT NOT NULL DEFAULT 'auth_bootstrap',
ADD COLUMN IF NOT EXISTS requested_by_user_id TEXT NULL;

UPDATE auth_wallet_challenges
SET purpose = COALESCE(NULLIF(TRIM(purpose), ''), 'auth_bootstrap')
WHERE purpose IS NULL OR TRIM(purpose) = '';

CREATE INDEX IF NOT EXISTS idx_wallet_challenges_purpose
ON auth_wallet_challenges(purpose);

CREATE INDEX IF NOT EXISTS idx_wallet_challenges_requested_by_user_id
ON auth_wallet_challenges(requested_by_user_id)
WHERE requested_by_user_id IS NOT NULL;