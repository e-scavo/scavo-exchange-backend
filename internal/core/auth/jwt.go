package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService struct {
	secret []byte
	issuer string
	ttl    time.Duration
}

type Claims struct {
	UserID string `json:"uid"`
	Email  string `json:"email,omitempty"`
	jwt.RegisteredClaims
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
	now := time.Now().UTC()
	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			Subject:   userID,
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
