package httpx

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewRouter_DiagnosticsEndpointExposesObservabilityState(t *testing.T) {
	router := newVersioningTestRouter(t)
	req := httptest.NewRequest(http.MethodGet, "/diagnostics", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: got=%d want=%d body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var payload struct {
		OK            bool `json:"ok"`
		Observability struct {
			RequestCorrelation     bool `json:"request_correlation"`
			StructuredLogging      bool `json:"structured_logging"`
			ErrorContextEnrichment bool `json:"error_context_enrichment"`
			FlowTracing            bool `json:"flow_tracing"`
		} `json:"observability"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatalf("decode error: %v body=%s", err, rec.Body.String())
	}

	if !payload.OK {
		t.Fatalf("expected diagnostics ok=true")
	}
	if !payload.Observability.RequestCorrelation || !payload.Observability.StructuredLogging || !payload.Observability.ErrorContextEnrichment || !payload.Observability.FlowTracing {
		t.Fatalf("unexpected observability state: %+v", payload.Observability)
	}
}
