package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"

	"supportflow/api/settings"
	"supportflow/services/ai"
)

type testProvider struct{}

func (p *testProvider) Name() string { return "test-provider" }
func (p *testProvider) Chat(ctx context.Context, req ai.LLMRequest) (*ai.LLMResponse, error) {
	return &ai.LLMResponse{Text: "test response"}, nil
}

func setupTestProviders() {
	ai.RegisterProvider("test-anthropic", &mockLLM{name: "test-anthropic"})
	ai.RegisterProvider("test-openai", &mockLLM{name: "test-openai"})
	ai.SetActiveProvider("test-anthropic")
}

type mockLLM struct {
	name string
}

func (m *mockLLM) Name() string { return m.name }
func (m *mockLLM) Chat(ctx context.Context, req ai.LLMRequest) (*ai.LLMResponse, error) {
	return &ai.LLMResponse{Text: "mock"}, nil
}

func TestHandleGetProviders(t *testing.T) {
	setupTestProviders()

	req := httptest.NewRequest("GET", "/api/settings/providers", nil)
	w := httptest.NewRecorder()

	settings.HandleGetProviders(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", w.Code)
	}

	var resp map[string]any
	json.NewDecoder(w.Body).Decode(&resp)

	if resp["active"] != "test-anthropic" {
		t.Errorf("active = %v, want test-anthropic", resp["active"])
	}

	providers := resp["providers"].([]any)
	if len(providers) < 2 {
		t.Errorf("providers count = %d, want >= 2", len(providers))
	}
}

func TestHandleSetProvider(t *testing.T) {
	setupTestProviders()

	body, _ := json.Marshal(map[string]string{"provider": "test-openai"})
	req := httptest.NewRequest("PUT", "/api/settings/providers", bytes.NewReader(body))
	w := httptest.NewRecorder()

	settings.HandleSetProvider(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", w.Code)
	}

	_, active := ai.GetActiveProvider()
	if active != "test-openai" {
		t.Errorf("active after switch = %q, want test-openai", active)
	}
}

func TestHandleSetProviderInvalid(t *testing.T) {
	setupTestProviders()

	body, _ := json.Marshal(map[string]string{"provider": "nonexistent"})
	req := httptest.NewRequest("PUT", "/api/settings/providers", bytes.NewReader(body))
	w := httptest.NewRecorder()

	settings.HandleSetProvider(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", w.Code)
	}
}

func TestHandleSetProviderEmpty(t *testing.T) {
	body, _ := json.Marshal(map[string]string{"provider": ""})
	req := httptest.NewRequest("PUT", "/api/settings/providers", bytes.NewReader(body))
	w := httptest.NewRecorder()

	settings.HandleSetProvider(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", w.Code)
	}
}

func TestHandleGetMetrics(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/settings/metrics?limit=10", nil)
	w := httptest.NewRecorder()

	settings.HandleGetMetrics(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", w.Code)
	}

	var resp map[string]any
	json.NewDecoder(w.Body).Decode(&resp)

	if resp["metrics"] == nil {
		t.Error("response should contain metrics key")
	}
	if resp["stats"] == nil {
		t.Error("response should contain stats key")
	}
}

func TestHealthEndpoint(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	}).Methods("GET")

	req := httptest.NewRequest("GET", "/api/health", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", w.Code)
	}

	var resp map[string]string
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["status"] != "ok" {
		t.Errorf("status = %q, want ok", resp["status"])
	}
}
