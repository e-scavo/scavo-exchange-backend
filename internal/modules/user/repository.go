package user

import "context"

type Repository interface {
	UpsertDevUser(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
}
