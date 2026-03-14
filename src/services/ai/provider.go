package ai

import (
	"context"
	"sync"
	"time"
)

type LLMRequest struct {
	SystemPrompt string
	Messages     []LLMMessage
	Tools        []ToolDef
	MaxTokens    int
}

type LLMMessage struct {
	Role       string
	Content    string
	ToolCalls  []ToolCall
	ToolResult *ToolResultMsg
}

type ToolCall struct {
	ID     string
	Name   string
	Params string
}

type ToolResultMsg struct {
	ToolUseID string
	Content   string
}

type ToolDef struct {
	Name        string
	Description string
	Parameters  map[string]any
	Required    []string
}

type LLMResponse struct {
	Text      string
	ToolCalls []ToolCall
	Usage     TokenUsage
}

type TokenUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

type LLMProvider interface {
	Name() string
	Chat(ctx context.Context, req LLMRequest) (*LLMResponse, error)
}

type LLMMetrics struct {
	Provider     string     `json:"provider"`
	Model        string     `json:"model"`
	LatencyMs    int64      `json:"latency_ms"`
	InputTokens  int        `json:"input_tokens"`
	OutputTokens int        `json:"output_tokens"`
	ToolCalls    int        `json:"tool_calls"`
	Timestamp    time.Time  `json:"timestamp"`
	Error        string     `json:"error,omitempty"`
}

var (
	providers   = map[string]LLMProvider{}
	activeProvider string
	providerMu sync.RWMutex
	metricsLog []LLMMetrics
	metricsMu  sync.Mutex
)

func RegisterProvider(name string, p LLMProvider) {
	providerMu.Lock()
	defer providerMu.Unlock()
	providers[name] = p
	if activeProvider == "" {
		activeProvider = name
	}
}

func SetActiveProvider(name string) bool {
	providerMu.Lock()
	defer providerMu.Unlock()
	if _, ok := providers[name]; !ok {
		return false
	}
	activeProvider = name
	return true
}

func GetActiveProvider() (LLMProvider, string) {
	providerMu.RLock()
	defer providerMu.RUnlock()
	return providers[activeProvider], activeProvider
}

func GetProviderNames() []string {
	providerMu.RLock()
	defer providerMu.RUnlock()
	names := make([]string, 0, len(providers))
	for n := range providers {
		names = append(names, n)
	}
	return names
}

func LogMetrics(m LLMMetrics) {
	metricsMu.Lock()
	defer metricsMu.Unlock()
	metricsLog = append(metricsLog, m)
}

func GetMetrics(limit int) []LLMMetrics {
	metricsMu.Lock()
	defer metricsMu.Unlock()
	if limit <= 0 || limit > len(metricsLog) {
		limit = len(metricsLog)
	}
	start := len(metricsLog) - limit
	result := make([]LLMMetrics, limit)
	copy(result, metricsLog[start:])
	return result
}

func GetMetricsStats() map[string]any {
	metricsMu.Lock()
	defer metricsMu.Unlock()

	if len(metricsLog) == 0 {
		return map[string]any{"total_calls": 0}
	}

	byProvider := map[string][]int64{}
	providerTokens := map[string][2]int{}
	providerErrors := map[string]int{}
	totalTokensIn := 0
	totalTokensOut := 0
	totalToolCalls := 0
	errors := 0

	for _, m := range metricsLog {
		byProvider[m.Provider] = append(byProvider[m.Provider], m.LatencyMs)
		totalTokensIn += m.InputTokens
		totalTokensOut += m.OutputTokens
		totalToolCalls += m.ToolCalls
		t := providerTokens[m.Provider]
		t[0] += m.InputTokens
		t[1] += m.OutputTokens
		providerTokens[m.Provider] = t
		if m.Error != "" {
			errors++
			providerErrors[m.Provider]++
		}
	}

	providerStats := map[string]any{}
	for name, latencies := range byProvider {
		var sum int64
		var min, max int64
		min = latencies[0]
		for _, l := range latencies {
			sum += l
			if l < min {
				min = l
			}
			if l > max {
				max = l
			}
		}
		t := providerTokens[name]
		providerStats[name] = map[string]any{
			"calls":        len(latencies),
			"avg_ms":       sum / int64(len(latencies)),
			"min_ms":       min,
			"max_ms":       max,
			"total_tokens": t[0] + t[1],
			"errors":       providerErrors[name],
		}
	}

	return map[string]any{
		"total_calls":       len(metricsLog),
		"total_tokens_in":   totalTokensIn,
		"total_tokens_out":  totalTokensOut,
		"total_tool_calls":  totalToolCalls,
		"errors":            errors,
		"by_provider":       providerStats,
	}
}
