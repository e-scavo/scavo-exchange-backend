package logger

import "testing"

func TestAttrsWithRequestID_ReturnsStandardAttribute(t *testing.T) {
	attrs := AttrsWithRequestID("request-123")

	if len(attrs) != 2 {
		t.Fatalf("unexpected attrs length: got=%d want=2", len(attrs))
	}
	if attrs[0] != RequestIDKey {
		t.Fatalf("unexpected attr key: got=%v want=%s", attrs[0], RequestIDKey)
	}
	if attrs[1] != "request-123" {
		t.Fatalf("unexpected attr value: got=%v want=request-123", attrs[1])
	}
}

func TestAttrsWithRequestID_ReturnsNilForEmptyRequestID(t *testing.T) {
	if attrs := AttrsWithRequestID(""); attrs != nil {
		t.Fatalf("expected nil attrs for empty request id, got=%v", attrs)
	}
}

func TestWithRequestID_PreservesNilOrEmptyInputs(t *testing.T) {
	if got := WithRequestID(nil, "request-123"); got != nil {
		t.Fatalf("expected nil logger to remain nil")
	}

	log := New("test")
	if got := WithRequestID(log, ""); got != log {
		t.Fatalf("expected empty request id to leave logger unchanged")
	}
}

func TestWithRequestID_ReturnsDerivedLogger(t *testing.T) {
	log := New("test")
	got := WithRequestID(log, "request-123")

	if got == nil {
		t.Fatalf("expected derived logger")
	}
	if got == log {
		t.Fatalf("expected a derived logger instance")
	}
}
