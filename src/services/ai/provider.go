package ai

import (
	"context"
	"math/rand"
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
	InputTokens      int `json:"input_tokens"`
	OutputTokens     int `json:"output_tokens"`
	CacheWriteTokens int `json:"cache_write_tokens"`
	CacheReadTokens  int `json:"cache_read_tokens"`
}

type LLMProvider interface {
	Name() string
	Chat(ctx context.Context, req LLMRequest) (*LLMResponse, error)
}

type LLMMetrics struct {
	Provider         string    `json:"provider"`
	Model            string    `json:"model"`
	LatencyMs        int64     `json:"latency_ms"`
	InputTokens      int       `json:"input_tokens"`
	OutputTokens     int       `json:"output_tokens"`
	CacheWriteTokens int       `json:"cache_write_tokens"`
	CacheReadTokens  int       `json:"cache_read_tokens"`
	ToolCalls        int       `json:"tool_calls"`
	Timestamp        time.Time `json:"timestamp"`
	Error            string    `json:"error,omitempty"`
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

func SeedDemoMetrics() {
	now := time.Now()

	type providerProfile struct {
		name, model                          string
		latencyMin, latencyMax               int
		inputMin, inputMax                   int
		outputMin, outputMax                 int
		cacheWriteMin, cacheWriteMax         int
		cacheReadMin, cacheReadMax           int
		toolCallChance                       float64
	}

	profiles := []providerProfile{
		{"anthropic", "claude-haiku-4-5-20251001", 400, 1200, 800, 3000, 150, 600, 200, 1000, 600, 2500, 0.4},
		{"openai", "gpt-4o", 600, 1800, 1000, 3800, 180, 700, 0, 200, 500, 2800, 0.35},
	}

	for i := 0; i < 40; i++ {
		p := profiles[i%2]
		if i%5 == 0 {
			p = profiles[1]
		}

		latency := int64(p.latencyMin + rand.Intn(p.latencyMax-p.latencyMin))
		input := p.inputMin + rand.Intn(p.inputMax-p.inputMin)
		output := p.outputMin + rand.Intn(p.outputMax-p.outputMin)
		cacheWrite := p.cacheWriteMin + rand.Intn(p.cacheWriteMax-p.cacheWriteMin+1)
		cacheRead := p.cacheReadMin + rand.Intn(p.cacheReadMax-p.cacheReadMin+1)
		tools := 0
		if rand.Float64() < p.toolCallChance {
			tools = 1 + rand.Intn(3)
		}

		var errStr string
		if rand.Float64() < 0.03 {
			errStr = "rate_limit_exceeded"
			latency = int64(p.latencyMax + rand.Intn(500))
		}

		m := LLMMetrics{
			Provider:         p.name,
			Model:            p.model,
			LatencyMs:        latency,
			InputTokens:      input,
			OutputTokens:     output,
			CacheWriteTokens: cacheWrite,
			CacheReadTokens:  cacheRead,
			ToolCalls:        tools,
			Timestamp:        now.Add(-time.Duration(40-i) * 3 * time.Minute),
			Error:            errStr,
		}
		metricsLog = append(metricsLog, m)
	}
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
	providerCache := map[string][2]int{}
	providerErrors := map[string]int{}
	totalTokensIn := 0
	totalTokensOut := 0
	totalCacheWrite := 0
	totalCacheRead := 0
	totalToolCalls := 0
	errors := 0

	for _, m := range metricsLog {
		byProvider[m.Provider] = append(byProvider[m.Provider], m.LatencyMs)
		totalTokensIn += m.InputTokens
		totalTokensOut += m.OutputTokens
		totalCacheWrite += m.CacheWriteTokens
		totalCacheRead += m.CacheReadTokens
		totalToolCalls += m.ToolCalls
		t := providerTokens[m.Provider]
		t[0] += m.InputTokens
		t[1] += m.OutputTokens
		providerTokens[m.Provider] = t
		c := providerCache[m.Provider]
		c[0] += m.CacheWriteTokens
		c[1] += m.CacheReadTokens
		providerCache[m.Provider] = c
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
		c := providerCache[name]
		providerStats[name] = map[string]any{
			"calls":              len(latencies),
			"avg_ms":             sum / int64(len(latencies)),
			"min_ms":             min,
			"max_ms":             max,
			"total_tokens":       t[0] + t[1],
			"cache_write_tokens": c[0],
			"cache_read_tokens":  c[1],
			"errors":             providerErrors[name],
		}
	}

	return map[string]any{
		"total_calls":         len(metricsLog),
		"total_tokens_in":     totalTokensIn,
		"total_tokens_out":    totalTokensOut,
		"total_cache_write":   totalCacheWrite,
		"total_cache_read":    totalCacheRead,
		"total_tool_calls":    totalToolCalls,
		"errors":              errors,
		"by_provider":         providerStats,
	}
}
