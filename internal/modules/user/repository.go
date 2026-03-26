package user

import "context"

type Repository interface {
	UpsertDevUser(ctx context.Context, email string) (*User, error)
	UpsertWalletUser(ctx context.Context, id, email, displayName string) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
}
