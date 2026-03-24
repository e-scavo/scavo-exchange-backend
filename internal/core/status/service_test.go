package status

import (
	"context"
	"errors"
	"net/http"
	"testing"
)

func TestReadiness_OKWhenRequiredDependenciesAreUp(t *testing.T) {
	svc := New(
		"scavo-exchange-backend",
		"test",
		"dev",
		"",
		FuncChecker{
			NameValue:     "postgres",
			RequiredValue: true,
			Fn: func(ctx context.Context) error {
				return nil
			},
		},
	)

	code, payload := svc.Readiness(context.Background())

	if code != http.StatusOK {
		t.Fatalf("expected 200, got %d", code)
	}

	ok, _ := payload["ok"].(bool)
	if !ok {
		t.Fatalf("expected readiness ok=true, got payload=%v", payload)
	}
}

func TestReadiness_FailsWhenRequiredDependencyIsDown(t *testing.T) {
	svc := New(
		"scavo-exchange-backend",
		"test",
		"dev",
		"",
		FuncChecker{
			NameValue:     "postgres",
			RequiredValue: true,
			Fn: func(ctx context.Context) error {
				return errors.New("db down")
			},
		},
	)

	code, payload := svc.Readiness(context.Background())

	if code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", code)
	}

	ok, _ := payload["ok"].(bool)
	if ok {
		t.Fatalf("expected readiness ok=false, got payload=%v", payload)
	}
}

func TestReadiness_AllowsOptionalDependencyToBeDown(t *testing.T) {
	svc := New(
		"scavo-exchange-backend",
		"test",
		"dev",
		"",
		FuncChecker{
			NameValue:     "redis",
			RequiredValue: false,
			Fn: func(ctx context.Context) error {
				return errors.New("redis down")
			},
		},
	)

	code, payload := svc.Readiness(context.Background())

	if code != http.StatusOK {
		t.Fatalf("expected 200, got %d", code)
	}

	ok, _ := payload["ok"].(bool)
	if !ok {
		t.Fatalf("expected readiness ok=true, got payload=%v", payload)
	}
}
