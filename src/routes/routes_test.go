package routes

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestRegisterCreatesAllRoutes(t *testing.T) {
	r := mux.NewRouter()
	Register(r, context.Background())

	expected := []struct {
		method string
		path   string
	}{
		{"GET", "/api/health"},
		{"POST", "/api/chat"},
		{"GET", "/api/tickets"},
		{"GET", "/api/tickets/some-id"},
		{"PUT", "/api/tickets/some-id/status"},
		{"PUT", "/api/tickets/some-id/assign"},
		{"POST", "/api/tickets/actions/approve"},
		{"GET", "/api/analytics/overview"},
		{"GET", "/api/analytics/agents"},
		{"GET", "/api/settings/providers"},
		{"PUT", "/api/settings/providers"},
		{"GET", "/api/settings/metrics"},
	}

	for _, tc := range expected {
		req := httptest.NewRequest(tc.method, tc.path, nil)
		match := &mux.RouteMatch{}
		if !r.Match(req, match) {
			t.Errorf("route %s %s not registered", tc.method, tc.path)
		}
	}
}

func TestHealthEndpointViaRouter(t *testing.T) {
	r := mux.NewRouter()
	Register(r, context.Background())

	req := httptest.NewRequest("GET", "/api/health", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("health status = %d, want 200", w.Code)
	}
	if w.Body.String() != `{"status":"ok"}` {
		t.Errorf("health body = %q, want {\"status\":\"ok\"}", w.Body.String())
	}
}

func TestUnknownRouteReturns405or404(t *testing.T) {
	r := mux.NewRouter()
	Register(r, context.Background())

	req := httptest.NewRequest("DELETE", "/api/health", nil)
	match := &mux.RouteMatch{}
	if r.Match(req, match) && match.MatchErr == nil {
		t.Error("DELETE /api/health should not match any route")
	}
}
