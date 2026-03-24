package status

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/e-scavo/scavo-exchange-backend/internal/core/cache"
	"github.com/e-scavo/scavo-exchange-backend/internal/core/db"
)

type Checker interface {
	Name() string
	Required() bool
	Check(ctx context.Context) error
}

type FuncChecker struct {
	NameValue     string
	RequiredValue bool
	Fn            func(ctx context.Context) error
}

func (f FuncChecker) Name() string {
	return f.NameValue
}

func (f FuncChecker) Required() bool {
	return f.RequiredValue
}

func (f FuncChecker) Check(ctx context.Context) error {
	if f.Fn == nil {
		return nil
	}
	return f.Fn(ctx)
}

type DependencyResult struct {
	Name     string `json:"name"`
	Status   string `json:"status"`
	Required bool   `json:"required"`
	Error    string `json:"error,omitempty"`
}

type Service struct {
	appName  string
	env      string
	version  string
	commit   string
	checkers []Checker
}

func New(appName, env, version, commit string, checkers ...Checker) *Service {
	return &Service{
		appName:  appName,
		env:      env,
		version:  version,
		commit:   commit,
		checkers: checkers,
	}
}

func (s *Service) Health() map[string]any {
	return map[string]any{
		"ok":      true,
		"status":  "up",
		"service": s.appName,
		"env":     s.env,
		"version": s.version,
		"commit":  s.commit,
		"time":    time.Now().UTC().Format(time.RFC3339),
	}
}

func (s *Service) Readiness(ctx context.Context) (int, map[string]any) {
	results := make([]DependencyResult, 0, len(s.checkers))
	ready := true

	for _, checker := range s.checkers {
		err := checker.Check(ctx)

		switch {
		case err == nil:
			results = append(results, DependencyResult{
				Name:     checker.Name(),
				Status:   "up",
				Required: checker.Required(),
			})

		case errors.Is(err, db.ErrNotConfigured), errors.Is(err, cache.ErrNotConfigured):
			results = append(results, DependencyResult{
				Name:     checker.Name(),
				Status:   "not_configured",
				Required: checker.Required(),
			})

			if checker.Required() {
				ready = false
			}

		default:
			results = append(results, DependencyResult{
				Name:     checker.Name(),
				Status:   "down",
				Required: checker.Required(),
				Error:    err.Error(),
			})

			if checker.Required() {
				ready = false
			}
		}
	}

	statusCode := http.StatusOK
	if !ready {
		statusCode = http.StatusServiceUnavailable
	}

	return statusCode, map[string]any{
		"ok":           ready,
		"status":       ternary(ready, "ready", "not_ready"),
		"service":      s.appName,
		"env":          s.env,
		"version":      s.version,
		"commit":       s.commit,
		"dependencies": results,
		"time":         time.Now().UTC().Format(time.RFC3339),
	}
}

func ternary[T any](cond bool, a, b T) T {
	if cond {
		return a
	}
	return b
}
