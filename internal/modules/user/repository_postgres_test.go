package user

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/e-scavo/scavo-exchange-backend/internal/core/logger"
)

func TestPostgresRepository_UpsertDevUser(t *testing.T) {
	dsn := os.Getenv("SCAVO_TEST_POSTGRES_URL")
	if dsn == "" {
		t.Skip("SCAVO_TEST_POSTGRES_URL not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("pgxpool.New error: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		t.Fatalf("postgres ping error: %v", err)
	}

	lg := logger.New("test")
	repo := NewPostgresRepository(pool, lg)

	email := "repo_test_user@example.com"

	u1, err := repo.UpsertDevUser(ctx, email)
	if err != nil {
		t.Fatalf("first upsert error: %v", err)
	}
	if u1 == nil {
		t.Fatal("expected user from first upsert")
	}
	if u1.Email != email {
		t.Fatalf("unexpected email after first upsert: %q", u1.Email)
	}
	if u1.LastLoginAt == nil {
		t.Fatal("expected LastLoginAt on first upsert")
	}

	time.Sleep(1100 * time.Millisecond)

	u2, err := repo.UpsertDevUser(ctx, email)
	if err != nil {
		t.Fatalf("second upsert error: %v", err)
	}
	if u2 == nil {
		t.Fatal("expected user from second upsert")
	}
	if u2.ID != u1.ID {
		t.Fatalf("expected same user id, got %q vs %q", u2.ID, u1.ID)
	}
	if u2.LastLoginAt == nil {
		t.Fatal("expected LastLoginAt on second upsert")
	}
	if !u2.LastLoginAt.After(*u1.LastLoginAt) && !u2.LastLoginAt.Equal(*u1.LastLoginAt) {
		t.Fatalf("expected second login timestamp >= first login timestamp")
	}
}

func TestPostgresRepository_UpdateDisplayName(t *testing.T) {
	dsn := os.Getenv("SCAVO_TEST_POSTGRES_URL")
	if dsn == "" {
		t.Skip("SCAVO_TEST_POSTGRES_URL not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("pgxpool.New error: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		t.Fatalf("postgres ping error: %v", err)
	}

	lg := logger.New("test")
	repo := NewPostgresRepository(pool, lg)

	u, err := repo.UpsertDevUser(ctx, "repo_update_display_name@example.com")
	if err != nil {
		t.Fatalf("upsert user error: %v", err)
	}
	if u == nil {
		t.Fatal("expected user from upsert")
	}

	time.Sleep(1100 * time.Millisecond)

	updated, err := repo.UpdateDisplayName(ctx, u.ID, "SCAVO Profile")
	if err != nil {
		t.Fatalf("update display name error: %v", err)
	}
	if updated == nil {
		t.Fatal("expected user from update")
	}
	if updated.DisplayName != "SCAVO Profile" {
		t.Fatalf("unexpected display name: %q", updated.DisplayName)
	}
	if !updated.UpdatedAt.After(u.UpdatedAt) && !updated.UpdatedAt.Equal(u.UpdatedAt) {
		t.Fatalf("expected updated timestamp >= previous updated timestamp")
	}
}
