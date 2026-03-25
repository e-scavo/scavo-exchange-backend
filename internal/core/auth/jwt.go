package auth

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService struct {
	secret []byte
	issuer string
	ttl    time.Duration
}

type Claims struct {
	UserID        string `json:"uid"`
	Email         string `json:"email,omitempty"`
	WalletID      string `json:"wallet_id,omitempty"`
	WalletAddress string `json:"wallet_address,omitempty"`
	AuthMethod    string `json:"auth_method,omitempty"`
	Chain         string `json:"chain,omitempty"`
	jwt.RegisteredClaims
}

type MintOptions struct {
	UserID        string
	Email         string
	WalletID      string
	WalletAddress string
	AuthMethod    string
	Chain         string
	Subject       string
}

func NewTokenService(secret string, issuer string, ttl time.Duration) (*TokenService, error) {
	if len(secret) < 24 {
		return nil, errors.New("SCAVO_JWT_SECRET must be at least 24 chars")
	}
	if issuer == "" {
		issuer = "scavo-exchange-backend"
	}
	if ttl <= 0 {
		ttl = 24 * time.Hour
	}
	return &TokenService{secret: []byte(secret), issuer: issuer, ttl: ttl}, nil
}

func (s *TokenService) Mint(userID, email string) (string, error) {
	return s.MintWithOptions(MintOptions{
		UserID:     userID,
		Email:      email,
		AuthMethod: "password_dev",
	})
}

func (s *TokenService) MintWithOptions(opts MintOptions) (string, error) {
	now := time.Now().UTC()
	subject := strings.TrimSpace(opts.Subject)
	if subject == "" {
		subject = strings.TrimSpace(opts.UserID)
	}

	claims := Claims{
		UserID:        strings.TrimSpace(opts.UserID),
		Email:         strings.TrimSpace(opts.Email),
		WalletID:      strings.TrimSpace(opts.WalletID),
		WalletAddress: strings.TrimSpace(opts.WalletAddress),
		AuthMethod:    strings.TrimSpace(opts.AuthMethod),
		Chain:         strings.TrimSpace(opts.Chain),
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			Subject:   subject,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.ttl)),
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tok.SignedString(s.secret)
}

func (s *TokenService) Parse(tokenStr string) (*Claims, error) {
	tok, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, jwt.ErrSignatureInvalid
		}
		return s.secret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := tok.Claims.(*Claims)
	if !ok || !tok.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}
	return claims, nil
}
