package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"supportflow/core"
	"supportflow/core/constants"
	"supportflow/core/structs"
	"supportflow/db/postgre"
)

func Init() {
	if p := NewAnthropicProvider(); p != nil {
		RegisterProvider("anthropic", p)
		log.Println("Anthropic provider registered")
	}
	if p := NewOpenAIProvider(); p != nil {
		RegisterProvider("openai", p)
		log.Println("OpenAI provider registered")
	}

	defaultProvider := core.GetString("ai.provider", "anthropic")
	SetActiveProvider(defaultProvider)
	log.Printf("Active AI provider: %s", defaultProvider)
}

var tools = []ToolDef{
	{Name: "lookup_customer", Description: "Look up customer information including profile, plan, and account details", Parameters: map[string]any{"customer_id": map[string]any{"type": "string", "description": "Customer ID"}}},
	{Name: "lookup_billing", Description: "Look up billing history, payments, invoices for a customer", Parameters: map[string]any{"customer_id": map[string]any{"type": "string", "description": "Customer ID"}}},
	{Name: "refund", Description: "Process a refund for a customer payment", Parameters: map[string]any{"amount": map[string]any{"type": "number", "description": "Amount to refund"}, "reason": map[string]any{"type": "string", "description": "Reason for refund"}}, Required: []string{"amount", "reason"}},
	{Name: "change_plan", Description: "Change customer subscription plan", Parameters: map[string]any{"new_plan": map[string]any{"type": "string", "description": "New plan name (free, basic, premium)"}}, Required: []string{"new_plan"}},
	{Name: "reset_password", Description: "Send password reset link to customer email", Parameters: map[string]any{}},
	{Name: "escalate", Description: "Escalate ticket to senior support or manager", Parameters: map[string]any{"reason": map[string]any{"type": "string", "description": "Reason for escalation"}, "priority": map[string]any{"type": "string", "description": "Priority: low, medium, high, critical"}}, Required: []string{"reason"}},
	{Name: "send_email", Description: "Send email notification to customer", Parameters: map[string]any{"subject": map[string]any{"type": "string", "description": "Email subject"}, "body": map[string]any{"type": "string", "description": "Email body"}}, Required: []string{"subject", "body"}},
	{Name: "cancel_subscription", Description: "Cancel customer subscription", Parameters: map[string]any{}},
}

const systemPrompt = `You are SupportFlow AI, an intelligent customer support assistant.
You work alongside human support agents to resolve customer issues quickly and effectively.

Your capabilities:
- Understand customer intent and sentiment
- Look up customer and billing information
- Execute actions: refunds, plan changes, password resets, escalations, emails
- Provide clear reasoning for every decision

Guidelines:
- Always look up relevant customer/billing info before taking action
- For simple, clear-cut cases (obvious double charges, password resets): act with high confidence
- For ambiguous cases: explain options and recommend, but let the agent decide
- Always be empathetic and professional in responses to customers
- Include your reasoning chain in responses
- Respond in the same language as the customer's message

When you use tools, explain what you're doing and why.
After resolving, provide a brief summary of actions taken.`

func ProcessMessage(ctx context.Context, ticketID, customerMessage string) (*structs.ChatResponse, error) {
	provider, providerName := GetActiveProvider()
	if provider == nil {
		return nil, fmt.Errorf("no AI provider configured")
	}

	history, _ := postgre.GetMessagesByTicket(ctx, ticketID)

	var llmMessages []LLMMessage
	for _, m := range history {
		role := "user"
		if m.Role == constants.RoleAI || m.Role == constants.RoleAgent {
			role = "assistant"
		}
		llmMessages = append(llmMessages, LLMMessage{Role: role, Content: m.Content})
	}
	llmMessages = append(llmMessages, LLMMessage{Role: "user", Content: customerMessage})

	response := &structs.ChatResponse{TicketID: ticketID}
	var allActions []structs.Action

	for i := 0; i < 5; i++ {
		start := time.Now()

		resp, err := provider.Chat(ctx, LLMRequest{
			SystemPrompt: systemPrompt,
			Messages:     llmMessages,
			Tools:        tools,
			MaxTokens:    core.GetInt("anthropic.max_tokens", 4096),
		})

		latency := time.Since(start).Milliseconds()

		metrics := LLMMetrics{
			Provider:  providerName,
			Model:     core.GetString(providerName+".model", "unknown"),
			LatencyMs: latency,
			Timestamp: time.Now(),
		}

		if err != nil {
			metrics.Error = err.Error()
			LogMetrics(metrics)
			return nil, fmt.Errorf("LLM error (%s): %w", providerName, err)
		}

		metrics.InputTokens = resp.Usage.InputTokens
		metrics.OutputTokens = resp.Usage.OutputTokens
		metrics.ToolCalls = len(resp.ToolCalls)
		LogMetrics(metrics)

		log.Printf("[%s] latency=%dms tokens_in=%d tokens_out=%d tools=%d",
			providerName, latency, resp.Usage.InputTokens, resp.Usage.OutputTokens, len(resp.ToolCalls))

		if resp.Text != "" {
			if response.Message != "" {
				response.Message += "\n"
			}
			response.Message += resp.Text
		}

		if len(resp.ToolCalls) == 0 {
			break
		}

		assistantMsg := LLMMessage{
			Role:      "assistant",
			Content:   resp.Text,
			ToolCalls: resp.ToolCalls,
		}
		llmMessages = append(llmMessages, assistantMsg)

		for _, tc := range resp.ToolCalls {
			result := ExecuteTool(ctx, ticketID, tc.Name, tc.Params)
			resultJSON, _ := json.Marshal(result)

			action := structs.Action{
				TicketID:   ticketID,
				Type:       tc.Name,
				Params:     tc.Params,
				Status:     constants.ActionStatusExecuted,
				Result:     string(resultJSON),
				Confidence: 0.9,
			}
			postgre.CreateAction(ctx, &action)
			allActions = append(allActions, action)

			llmMessages = append(llmMessages, LLMMessage{
				Role: "user",
				ToolResult: &ToolResultMsg{
					ToolUseID: tc.ID,
					Content:   string(resultJSON),
				},
			})
		}
	}

	response.Actions = allActions
	if len(allActions) > 0 {
		response.AutoFixed = true
	}

	postgre.UpdateTicketSummary(ctx, ticketID, response.Message)

	return response, nil
}

func AnalyzeTicket(ctx context.Context, ticketID, message string) (*structs.AIAnalysis, error) {
	provider, providerName := GetActiveProvider()
	if provider == nil {
		return nil, fmt.Errorf("no AI provider configured")
	}

	prompt := fmt.Sprintf(`Analyze this customer support message and return a JSON object:
{
  "intent": "<billing_dispute|technical_issue|account_access|plan_change|cancellation|general_inquiry|complaint>",
  "sentiment": "<positive|neutral|negative|angry>",
  "urgency": "<low|medium|high>",
  "suggested_tools": ["<tool_names>"],
  "reasoning": "<brief explanation>",
  "confidence": <0.0-1.0>
}

Customer message: "%s"

Return ONLY the JSON, no extra text.`, message)

	start := time.Now()
	resp, err := provider.Chat(ctx, LLMRequest{
		Messages:  []LLMMessage{{Role: "user", Content: prompt}},
		MaxTokens: 1024,
	})
	latency := time.Since(start).Milliseconds()

	LogMetrics(LLMMetrics{
		Provider:     providerName,
		Model:        core.GetString(providerName+".model", "unknown"),
		LatencyMs:    latency,
		InputTokens:  resp.Usage.InputTokens,
		OutputTokens: resp.Usage.OutputTokens,
		Timestamp:    time.Now(),
	})

	if err != nil {
		return nil, err
	}

	var analysis structs.AIAnalysis
	if err := json.Unmarshal([]byte(resp.Text), &analysis); err != nil {
		analysis = structs.AIAnalysis{
			TicketID:   ticketID,
			Intent:     "general_inquiry",
			Sentiment:  "neutral",
			Urgency:    "medium",
			Reasoning:  resp.Text,
			Confidence: 0.5,
		}
	}
	analysis.TicketID = ticketID

	return &analysis, nil
}
