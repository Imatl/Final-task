package ai

import (
	"context"
	"testing"
	"time"
)

type mockProvider struct {
	name     string
	response *LLMResponse
	err      error
}

func (m *mockProvider) Name() string { return m.name }
func (m *mockProvider) Chat(ctx context.Context, req LLMRequest) (*LLMResponse, error) {
	return m.response, m.err
}

func resetProviders() {
	providerMu.Lock()
	providers = map[string]LLMProvider{}
	activeProvider = ""
	providerMu.Unlock()

	metricsMu.Lock()
	metricsLog = nil
	metricsMu.Unlock()
}

func TestRegisterProvider(t *testing.T) {
	resetProviders()

	p := &mockProvider{name: "test"}
	RegisterProvider("test", p)

	names := GetProviderNames()
	if len(names) != 1 || names[0] != "test" {
		t.Errorf("providers = %v, want [test]", names)
	}

	active, name := GetActiveProvider()
	if name != "test" || active == nil {
		t.Errorf("active = %q, want test", name)
	}
}

func TestSetActiveProvider(t *testing.T) {
	resetProviders()

	RegisterProvider("a", &mockProvider{name: "a"})
	RegisterProvider("b", &mockProvider{name: "b"})

	if !SetActiveProvider("b") {
		t.Error("SetActiveProvider(b) should return true")
	}
	_, name := GetActiveProvider()
	if name != "b" {
		t.Errorf("active = %q, want b", name)
	}

	if SetActiveProvider("nonexistent") {
		t.Error("SetActiveProvider(nonexistent) should return false")
	}
	_, name = GetActiveProvider()
	if name != "b" {
		t.Errorf("active should remain b, got %q", name)
	}
}

func TestGetProviderNames(t *testing.T) {
	resetProviders()

	RegisterProvider("alpha", &mockProvider{name: "alpha"})
	RegisterProvider("beta", &mockProvider{name: "beta"})

	names := GetProviderNames()
	if len(names) != 2 {
		t.Errorf("len(names) = %d, want 2", len(names))
	}

	nameSet := map[string]bool{}
	for _, n := range names {
		nameSet[n] = true
	}
	if !nameSet["alpha"] || !nameSet["beta"] {
		t.Errorf("names = %v, want alpha and beta", names)
	}
}

func TestLogAndGetMetrics(t *testing.T) {
	resetProviders()

	LogMetrics(LLMMetrics{Provider: "test", LatencyMs: 100, InputTokens: 50, OutputTokens: 20, Timestamp: time.Now()})
	LogMetrics(LLMMetrics{Provider: "test", LatencyMs: 200, InputTokens: 60, OutputTokens: 30, Timestamp: time.Now()})
	LogMetrics(LLMMetrics{Provider: "other", LatencyMs: 150, InputTokens: 40, OutputTokens: 10, Timestamp: time.Now()})

	all := GetMetrics(0)
	if len(all) != 3 {
		t.Errorf("GetMetrics(0) = %d entries, want 3", len(all))
	}

	last2 := GetMetrics(2)
	if len(last2) != 2 {
		t.Errorf("GetMetrics(2) = %d entries, want 2", len(last2))
	}
	if last2[0].LatencyMs != 200 {
		t.Errorf("last2[0].LatencyMs = %d, want 200", last2[0].LatencyMs)
	}
}

func TestGetMetricsStats(t *testing.T) {
	resetProviders()

	LogMetrics(LLMMetrics{Provider: "anthropic", LatencyMs: 100, InputTokens: 50, OutputTokens: 20, ToolCalls: 1, Timestamp: time.Now()})
	LogMetrics(LLMMetrics{Provider: "anthropic", LatencyMs: 200, InputTokens: 60, OutputTokens: 30, ToolCalls: 2, Timestamp: time.Now()})
	LogMetrics(LLMMetrics{Provider: "openai", LatencyMs: 150, InputTokens: 40, OutputTokens: 10, Error: "timeout", Timestamp: time.Now()})

	stats := GetMetricsStats()

	if stats["total_calls"] != 3 {
		t.Errorf("total_calls = %v, want 3", stats["total_calls"])
	}
	if stats["errors"] != 1 {
		t.Errorf("errors = %v, want 1", stats["errors"])
	}
	if stats["total_tokens_in"] != 150 {
		t.Errorf("total_tokens_in = %v, want 150", stats["total_tokens_in"])
	}

	byProvider := stats["by_provider"].(map[string]any)
	anthropicStats := byProvider["anthropic"].(map[string]any)
	if anthropicStats["calls"] != 2 {
		t.Errorf("anthropic calls = %v, want 2", anthropicStats["calls"])
	}
	if anthropicStats["min_ms"] != int64(100) {
		t.Errorf("anthropic min_ms = %v, want 100", anthropicStats["min_ms"])
	}
	if anthropicStats["max_ms"] != int64(200) {
		t.Errorf("anthropic max_ms = %v, want 200", anthropicStats["max_ms"])
	}
}

func TestGetMetricsStatsEmpty(t *testing.T) {
	resetProviders()

	stats := GetMetricsStats()
	if stats["total_calls"] != 0 {
		t.Errorf("empty stats total_calls = %v, want 0", stats["total_calls"])
	}
}

func TestFirstRegisteredBecomesActive(t *testing.T) {
	resetProviders()

	RegisterProvider("first", &mockProvider{name: "first"})
	_, name := GetActiveProvider()
	if name != "first" {
		t.Errorf("first registered should be active, got %q", name)
	}

	RegisterProvider("second", &mockProvider{name: "second"})
	_, name = GetActiveProvider()
	if name != "first" {
		t.Errorf("active should remain first, got %q", name)
	}
}
