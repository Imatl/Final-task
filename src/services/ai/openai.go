package ai

import (
	"context"
	"fmt"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/shared"

	"supportflow/core"
)

type OpenAIProvider struct {
	client *openai.Client
	model  string
}

func NewOpenAIProvider() *OpenAIProvider {
	apiKey := core.GetString("openai.api_key", "")
	if apiKey == "" {
		return nil
	}
	c := openai.NewClient(option.WithAPIKey(apiKey))
	return &OpenAIProvider{
		client: &c,
		model:  core.GetString("openai.model", "gpt-4o"),
	}
}

func (o *OpenAIProvider) Name() string { return "openai" }

func (o *OpenAIProvider) Chat(ctx context.Context, req LLMRequest) (*LLMResponse, error) {
	messages := make([]openai.ChatCompletionMessageParamUnion, 0, len(req.Messages)+1)

	if req.SystemPrompt != "" {
		messages = append(messages, openai.SystemMessage(req.SystemPrompt))
	}

	for _, m := range req.Messages {
		switch {
		case m.ToolResult != nil:
			messages = append(messages, openai.ToolMessage(m.ToolResult.Content, m.ToolResult.ToolUseID))
		case len(m.ToolCalls) > 0:
			calls := make([]openai.ChatCompletionMessageToolCallParam, 0, len(m.ToolCalls))
			for _, tc := range m.ToolCalls {
				calls = append(calls, openai.ChatCompletionMessageToolCallParam{
					ID: tc.ID,
					Function: openai.ChatCompletionMessageToolCallFunctionParam{
						Name:      tc.Name,
						Arguments: tc.Params,
					},
				})
			}
			messages = append(messages, openai.ChatCompletionMessageParamUnion{
				OfAssistant: &openai.ChatCompletionAssistantMessageParam{
					Content:   openai.ChatCompletionAssistantMessageParamContentUnion{OfString: openai.String(m.Content)},
					ToolCalls: calls,
				},
			})
		default:
			if m.Role == "assistant" {
				messages = append(messages, openai.AssistantMessage(m.Content))
			} else {
				messages = append(messages, openai.UserMessage(m.Content))
			}
		}
	}

	tools := make([]openai.ChatCompletionToolParam, 0, len(req.Tools))
	for _, t := range req.Tools {
		props := t.Parameters
		if props == nil {
			props = map[string]any{}
		}
		schema := shared.FunctionParameters{
			"type":       "object",
			"properties": props,
		}
		if len(t.Required) > 0 {
			schema["required"] = t.Required
		}

		tools = append(tools, openai.ChatCompletionToolParam{
			Function: openai.FunctionDefinitionParam{
				Name:        t.Name,
				Description: openai.String(t.Description),
				Parameters:  schema,
			},
		})
	}

	params := openai.ChatCompletionNewParams{
		Model:    openai.ChatModel(o.model),
		Messages: messages,
	}
	if len(tools) > 0 {
		params.Tools = tools
	}

	resp, err := o.client.Chat.Completions.New(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("openai: %w", err)
	}

	result := &LLMResponse{
		Usage: TokenUsage{
			InputTokens:  int(resp.Usage.PromptTokens),
			OutputTokens: int(resp.Usage.CompletionTokens),
		},
	}

	if len(resp.Choices) > 0 {
		choice := resp.Choices[0]
		result.Text = choice.Message.Content
		for _, tc := range choice.Message.ToolCalls {
			result.ToolCalls = append(result.ToolCalls, ToolCall{
				ID:     tc.ID,
				Name:   tc.Function.Name,
				Params: tc.Function.Arguments,
			})
		}
	}

	return result, nil
}

