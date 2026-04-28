package readmodels

import "time"

// UserReadModel is the explicit output projection for user data leaving the
// user domain boundary. It must remain output-only and must not be reused as a
// write/input payload.
type UserReadModel struct {
	ID          string     `json:"id"`
	Email       string     `json:"email"`
	DisplayName string     `json:"display_name"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
}
