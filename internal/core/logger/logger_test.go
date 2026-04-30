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

func TestAttrsWithFlowEvent_IncludesEventAndRequestID(t *testing.T) {
	attrs := AttrsWithFlowEvent("request-123", "http_request_start", "method", "GET")

	want := []any{RequestIDKey, "request-123", FlowEventKey, "http_request_start", "method", "GET"}
	if len(attrs) != len(want) {
		t.Fatalf("unexpected attrs length: got=%d want=%d", len(attrs), len(want))
	}
	for i := range want {
		if attrs[i] != want[i] {
			t.Fatalf("unexpected attr at %d: got=%v want=%v", i, attrs[i], want[i])
		}
	}
}

func TestAttrsWithFlowEvent_AllowsMissingRequestID(t *testing.T) {
	attrs := AttrsWithFlowEvent("", "application_start")

	if len(attrs) != 2 {
		t.Fatalf("unexpected attrs length: got=%d want=2", len(attrs))
	}
	if attrs[0] != FlowEventKey || attrs[1] != "application_start" {
		t.Fatalf("unexpected attrs: got=%v", attrs)
	}
}
