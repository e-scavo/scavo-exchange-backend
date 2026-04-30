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

func TestDiagnostics_ReportsObservabilityCapabilities(t *testing.T) {
	svc := New("scavo-exchange-backend", "test", "dev", "abc123")

	payload := svc.Diagnostics()

	if ok, _ := payload["ok"].(bool); !ok {
		t.Fatalf("expected diagnostics ok=true, got payload=%v", payload)
	}
	if payload["service"] != "scavo-exchange-backend" {
		t.Fatalf("unexpected service: got=%v", payload["service"])
	}
	observability, ok := payload["observability"].(map[string]any)
	if !ok {
		t.Fatalf("expected observability map, got=%T", payload["observability"])
	}

	for _, key := range []string{
		"request_correlation",
		"structured_logging",
		"error_context_enrichment",
		"flow_tracing",
	} {
		if enabled, _ := observability[key].(bool); !enabled {
			t.Fatalf("expected observability %s=true, got payload=%v", key, observability)
		}
	}
}
