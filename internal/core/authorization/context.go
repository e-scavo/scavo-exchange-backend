package authorization

import "context"

type contextKey string

const SubjectContextKey contextKey = "authorization_subject"

func WithSubject(ctx context.Context, subject AuthorizationSubject) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	normalized := subject.Normalized()
	if normalized.UserID == "" && len(normalized.Roles) == 0 {
		return ctx
	}

	return context.WithValue(ctx, SubjectContextKey, &normalized)
}

func SubjectFromContext(ctx context.Context) (*AuthorizationSubject, bool) {
	if ctx == nil {
		return nil, false
	}

	subject, ok := ctx.Value(SubjectContextKey).(*AuthorizationSubject)
	if !ok || subject == nil {
		return nil, false
	}

	normalized := subject.Normalized()
	return &normalized, true
}
