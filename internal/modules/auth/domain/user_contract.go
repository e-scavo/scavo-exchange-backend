package domain

import (
	"context"

	usermod "github.com/e-scavo/scavo-exchange-backend/internal/modules/user"
)

// =====================================================
// USER CONTRACT (CROSS-MODULE)
// =====================================================
//
// This interface defines the minimal contract required by
// the auth module to interact with user module.
//
// IMPORTANT:
// - This is NOT a mirror of user.Service
// - Only required methods are exposed
// - Implementation is provided by user module
//

type UserProvider interface {
	GetByID(ctx context.Context, userID, email string) (*usermod.User, error)
	ResolveOrCreateDevUser(ctx context.Context, email string) (*usermod.User, error)
}
