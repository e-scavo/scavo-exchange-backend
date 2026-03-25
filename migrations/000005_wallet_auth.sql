-- Wallet Challenges
CREATE TABLE IF NOT EXISTS auth_wallet_challenges (
    id UUID PRIMARY KEY,
    address TEXT NOT NULL,
    chain TEXT NOT NULL,
    nonce TEXT NOT NULL,
    message TEXT NOT NULL,
    issued_at TIMESTAMP NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    used_at TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_wallet_challenges_address ON auth_wallet_challenges(address);
CREATE INDEX IF NOT EXISTS idx_wallet_challenges_expires_at ON auth_wallet_challenges(expires_at);
CREATE INDEX IF NOT EXISTS idx_wallet_challenges_used_at ON auth_wallet_challenges(used_at);

-- Wallet Identities
CREATE TABLE IF NOT EXISTS auth_wallet_identities (
    id UUID PRIMARY KEY,
    address TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);