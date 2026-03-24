package user

import "context"

type Repository interface {
	UpsertDevUser(ctx context.Context, email string) (*User, error)
}
