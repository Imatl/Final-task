package ai

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"

	"supportflow/core"
)

type AnthropicProvider struct {
	client anthropic.Client
	model  string
}

func NewAnthropicProvider() *AnthropicProvider {
	apiKey := core.GetString("anthropic.api_key", "")
	if apiKey == "" {
		return nil
	}
	opts := []option.RequestOption{option.WithAPIKey(apiKey)}
	baseURL := core.GetString("anthropic.base_url", "")
	if baseURL != "" {
		opts = append(opts, option.WithBaseURL(baseURL))
	}
	return &AnthropicProvider{
		client: anthropic.NewClient(opts...),
		model:  core.GetString("anthropic.model", "claude-sonnet-4-5-20250514"),
	}
}

func (a *AnthropicProvider) Name() string { return "anthropic" }

func (a *AnthropicProvider) Chat(ctx context.Context, req LLMRequest) (*LLMResponse, error) {
	messages := make([]anthropic.MessageParam, 0, len(req.Messages))

	for _, m := range req.Messages {
		switch {
		case m.ToolResult != nil:
			messages = append(messages, anthropic.MessageParam{
				Role: anthropic.MessageParamRoleUser,
				Content: []anthropic.ContentBlockParamUnion{
					anthropic.NewToolResultBlock(m.ToolResult.ToolUseID, m.ToolResult.Content, false),
				},
			})
		case len(m.ToolCalls) > 0:
			var blocks []anthropic.ContentBlockParamUnion
			if m.Content != "" {
				blocks = append(blocks, anthropic.NewTextBlock(m.Content))
			}
			for _, tc := range m.ToolCalls {
				blocks = append(blocks, anthropic.ContentBlockParamUnion{
					OfToolUse: &anthropic.ToolUseBlockParam{
						ID:    tc.ID,
						Name:  tc.Name,
						Input: json.RawMessage(tc.Params),
					},
				})
			}
			messages = append(messages, anthropic.MessageParam{
				Role:    anthropic.MessageParamRoleAssistant,
				Content: blocks,
			})
		default:
			role := anthropic.MessageParamRoleUser
			if m.Role == "assistant" {
				role = anthropic.MessageParamRoleAssistant
			}
			messages = append(messages, anthropic.MessageParam{
				Role: role,
				Content: []anthropic.ContentBlockParamUnion{
					anthropic.NewTextBlock(m.Content),
				},
			})
		}
	}

	tools := make([]anthropic.ToolUnionParam, 0, len(req.Tools))
	for _, t := range req.Tools {
		tools = append(tools, anthropic.ToolUnionParam{
			OfTool: &anthropic.ToolParam{
				Name:        t.Name,
				Description: anthropic.String(t.Description),
				InputSchema: anthropic.ToolInputSchemaParam{
					Properties: t.Parameters,
					Required:   t.Required,
				},
			},
		})
	}

	maxTokens := int64(req.MaxTokens)
	if maxTokens == 0 {
		maxTokens = 4096
	}

	params := anthropic.MessageNewParams{
		Model:     anthropic.Model(a.model),
		MaxTokens: maxTokens,
		Messages:  messages,
		Tools:     tools,
	}

	if req.SystemPrompt != "" {
		params.System = []anthropic.TextBlockParam{{Text: req.SystemPrompt}}
	}

	resp, err := a.client.Messages.New(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("anthropic: %w", err)
	}

	result := &LLMResponse{
		Usage: TokenUsage{
			InputTokens:      int(resp.Usage.InputTokens),
			OutputTokens:     int(resp.Usage.OutputTokens),
			CacheWriteTokens: int(resp.Usage.CacheCreationInputTokens),
			CacheReadTokens:  int(resp.Usage.CacheReadInputTokens),
		},
	}

	for _, block := range resp.Content {
		if block.Type == "text" {
			if result.Text != "" {
				result.Text += "\n"
			}
			result.Text += block.Text
		}
		if block.Type == "tool_use" {
			result.ToolCalls = append(result.ToolCalls, ToolCall{
				ID:     block.ID,
				Name:   block.Name,
				Params: string(block.Input),
			})
		}
	}

	return result, nil
}
