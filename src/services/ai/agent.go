package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
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

	SeedDemoMetrics()
	log.Println("Demo metrics seeded (40 records)")
}

var tools = []ToolDef{
	{Name: "lookup_customer", Description: "Look up current customer's profile, plan, and account details. Customer is automatically identified from the ticket.", Parameters: map[string]any{}},
	{Name: "lookup_billing", Description: "Look up current customer's billing history, payments, and invoices. Customer is automatically identified from the ticket.", Parameters: map[string]any{}},
	{Name: "refund", Description: "Process a refund for a customer payment", Parameters: map[string]any{"amount": map[string]any{"type": "number", "description": "Amount to refund"}, "reason": map[string]any{"type": "string", "description": "Reason for refund"}}, Required: []string{"amount", "reason"}},
	{Name: "change_plan", Description: "Change customer subscription plan", Parameters: map[string]any{"new_plan": map[string]any{"type": "string", "description": "New plan name (free, basic, premium)"}}, Required: []string{"new_plan"}},
	{Name: "reset_password", Description: "Send password reset link to customer email", Parameters: map[string]any{}},
	{Name: "escalate", Description: "Escalate ticket to senior support or manager", Parameters: map[string]any{"reason": map[string]any{"type": "string", "description": "Reason for escalation"}, "priority": map[string]any{"type": "string", "description": "Priority: low, medium, high, critical"}}, Required: []string{"reason"}},
	{Name: "send_email", Description: "Send email notification to customer", Parameters: map[string]any{"subject": map[string]any{"type": "string", "description": "Email subject"}, "body": map[string]any{"type": "string", "description": "Email body"}}, Required: []string{"subject", "body"}},
	{Name: "cancel_subscription", Description: "Cancel customer subscription", Parameters: map[string]any{}},
}

const systemPromptTpl = `You are Kairon, an intelligent customer support assistant.
You communicate directly with customers to resolve their issues quickly.

IMPORTANT: You already have access to the customer's account through the ticket system.
Do NOT ask the customer for their ID, email, or account details — use the lookup_customer tool to get this information when needed.

Your capabilities:
- Look up customer information automatically
- Execute actions: refunds, plan changes, password resets, escalations, emails
- Provide clear, helpful responses

Guidelines:
- For greetings or generic messages (like "hi", "hello", "start", "/start"): just greet the customer warmly and ask how you can help. Do NOT call any tools for greetings
- When the customer describes a specific problem or asks a question: call lookup_customer FIRST, then help resolve it
- ONLY report facts returned by tools. NEVER invent, assume, or hallucinate data not present in tool results
- NEVER include raw customer data (name, email, phone, plan, dates) in your response. Use it internally to perform actions, but do not display it to the customer
- For ambiguous cases: explain options and recommend a course of action
- Be empathetic, professional, and VERY concise — 2-3 sentences max. No numbered lists or bullet points unless necessary
- Never ask the customer for information you can look up yourself
- Do not use emojis
- Do NOT generate any intermediate thinking or status messages like "let me check" — just call the tools silently and respond with the result
%s
After taking action, briefly confirm what was done.`

func buildSystemPrompt(lang string) string {
	langRule := ""
	switch lang {
	case "ru", "uk":
		langRule = "- ALWAYS respond in Ukrainian regardless of the customer's language"
	default:
		if lang != "" {
			langRule = "- ALWAYS respond in English regardless of the customer's language"
		} else {
			langRule = "- Respond in the same language as the customer's message"
		}
	}
	return fmt.Sprintf(systemPromptTpl, langRule)
}

func ProcessMessage(ctx context.Context, ticketID, customerMessage, lang string) (*structs.ChatResponse, error) {
	provider, providerName := GetActiveProvider()
	if provider == nil {
		return nil, fmt.Errorf("no AI provider configured")
	}

	history, err := postgre.GetMessagesByTicket(ctx, ticketID)
	if err != nil {
		log.Printf("[ai] get message history error: %v", err)
	}

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

	systemPrompt := buildSystemPrompt(lang)
	ticket, err := postgre.GetTicket(ctx, ticketID)
	if err == nil && ticket.Company != nil && *ticket.Company != "" {
		entries, kbErr := postgre.ListKBEntries(ctx, *ticket.Company)
		if kbErr == nil && len(entries) > 0 {
			kbContext := "\n\nCompany Knowledge Base (use this to answer customer questions):\n"
			for _, e := range entries {
				kbContext += fmt.Sprintf("Q: %s\nA: %s\n\n", e.Question, e.Answer)
			}
			systemPrompt += kbContext
		}
	}

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
		metrics.CacheWriteTokens = resp.Usage.CacheWriteTokens
		metrics.CacheReadTokens = resp.Usage.CacheReadTokens
		metrics.ToolCalls = len(resp.ToolCalls)
		LogMetrics(metrics)

		log.Printf("[%s] latency=%dms tokens_in=%d tokens_out=%d cache_w=%d cache_r=%d tools=%d",
			providerName, latency, resp.Usage.InputTokens, resp.Usage.OutputTokens, resp.Usage.CacheWriteTokens, resp.Usage.CacheReadTokens, len(resp.ToolCalls))

		if len(resp.ToolCalls) == 0 {
			response.Message = resp.Text
			break
		}

		assistantMsg := LLMMessage{
			Role:      "assistant",
			Content:   resp.Text,
			ToolCalls: resp.ToolCalls,
		}
		llmMessages = append(llmMessages, assistantMsg)

		for _, tc := range resp.ToolCalls {
			requiresApproval := tc.Name == "refund" || tc.Name == "cancel_subscription"

			var result ToolResult
			var resultJSON []byte

			if requiresApproval {
				result = ToolResult{
					Success: false,
					Message: "Action requires agent approval before execution",
				}
				resultJSON, _ = json.Marshal(result)
				resultStr := string(resultJSON)
				action := structs.Action{
					TicketID:   ticketID,
					Type:       tc.Name,
					Params:     tc.Params,
					Status:     constants.ActionStatusPending,
					Result:     &resultStr,
					Confidence: 0.9,
				}
				if err := postgre.CreateAction(ctx, &action); err != nil {
					log.Printf("[ai] save pending action error: %v", err)
				}
				allActions = append(allActions, action)
				log.Printf("[ai] action %s requires approval for ticket %s", tc.Name, ticketID)
			} else {
				result = ExecuteTool(ctx, ticketID, tc.Name, tc.Params)
				resultJSON, _ = json.Marshal(result)
				resultStr := string(resultJSON)
				action := structs.Action{
					TicketID:   ticketID,
					Type:       tc.Name,
					Params:     tc.Params,
					Status:     constants.ActionStatusExecuted,
					Result:     &resultStr,
					Confidence: 0.9,
				}
				if err := postgre.CreateAction(ctx, &action); err != nil {
					log.Printf("[ai] save action error: %v", err)
				}
				allActions = append(allActions, action)
			}

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

		allSucceeded := true
		for _, a := range allActions {
			var result ToolResult
			if a.Result != nil {
				if err := json.Unmarshal([]byte(*a.Result), &result); err == nil && !result.Success {
					allSucceeded = false
					break
				}
			}
		}
		if allSucceeded {
			if err := postgre.UpdateTicketStatus(ctx, ticketID, constants.TicketStatusResolved); err != nil {
				log.Printf("[ai] auto-resolve ticket error: %v", err)
			} else {
				log.Printf("[ai] ticket %s auto-resolved after successful actions", ticketID)
				go GenerateTicketSummary(context.Background(), ticketID)
			}
		}
	}

	if err := postgre.UpdateTicketSummary(ctx, ticketID, response.Message); err != nil {
		log.Printf("[ai] update ticket summary error: %v", err)
	}

	return response, nil
}

const suggestPrompt = `You are an AI assistant helping a human support agent draft replies to customers.
Based on the conversation history and customer information, draft a professional reply that the agent can review, edit, and send.

Guidelines:
- Write as if YOU are the support agent (not an AI)
- Be empathetic, professional, and concise
- If the customer has a clear issue, suggest a concrete resolution
- Respond in the same language as the customer's last message
- Do not use emojis
- Do not mention that this is an AI-generated draft`

func GenerateSuggestion(ctx context.Context, ticketID string) (string, error) {
	provider, providerName := GetActiveProvider()
	if provider == nil {
		return "", fmt.Errorf("no AI provider configured")
	}

	history, err := postgre.GetMessagesByTicket(ctx, ticketID)
	if err != nil {
		return "", fmt.Errorf("get message history: %w", err)
	}

	var llmMessages []LLMMessage
	for _, m := range history {
		role := "user"
		if m.Role == constants.RoleAI || m.Role == constants.RoleAgent {
			role = "assistant"
		}
		llmMessages = append(llmMessages, LLMMessage{Role: role, Content: m.Content})
	}

	llmMessages = append(llmMessages, LLMMessage{
		Role:    "user",
		Content: "Based on the conversation above, draft a reply for me as the support agent. Write ONLY the reply text, nothing else.",
	})

	start := time.Now()
	resp, err := provider.Chat(ctx, LLMRequest{
		SystemPrompt: suggestPrompt,
		Messages:     llmMessages,
		MaxTokens:    core.GetInt("anthropic.max_tokens", 4096),
	})
	latency := time.Since(start).Milliseconds()

	LogMetrics(LLMMetrics{
		Provider:         providerName,
		Model:            core.GetString(providerName+".model", "unknown"),
		LatencyMs:        latency,
		InputTokens:      resp.Usage.InputTokens,
		OutputTokens:     resp.Usage.OutputTokens,
		CacheWriteTokens: resp.Usage.CacheWriteTokens,
		CacheReadTokens:  resp.Usage.CacheReadTokens,
		Timestamp:        time.Now(),
	})

	if err != nil {
		return "", fmt.Errorf("LLM error (%s): %w", providerName, err)
	}

	log.Printf("[ai] suggestion generated for ticket %s latency=%dms", ticketID, latency)
	return resp.Text, nil
}

func GenerateTicketSummary(ctx context.Context, ticketID string) {
	provider, providerName := GetActiveProvider()
	if provider == nil {
		log.Printf("[ai] no provider for summary generation")
		return
	}

	history, err := postgre.GetMessagesByTicket(ctx, ticketID)
	if err != nil {
		log.Printf("[ai] summary: get messages error: %v", err)
		return
	}

	if len(history) == 0 {
		return
	}

	var conversation string
	for _, m := range history {
		conversation += fmt.Sprintf("[%s]: %s\n", m.Role, m.Content)
	}

	prompt := fmt.Sprintf(`Summarize this support conversation in 2-3 sentences. Include: what the customer wanted, what actions were taken, and the outcome.

Conversation:
%s

Return ONLY the summary text, nothing else.`, conversation)

	start := time.Now()
	resp, err := provider.Chat(ctx, LLMRequest{
		Messages:  []LLMMessage{{Role: "user", Content: prompt}},
		MaxTokens: 512,
	})
	latency := time.Since(start).Milliseconds()

	LogMetrics(LLMMetrics{
		Provider:         providerName,
		Model:            core.GetString(providerName+".model", "unknown"),
		LatencyMs:        latency,
		InputTokens:      resp.Usage.InputTokens,
		OutputTokens:     resp.Usage.OutputTokens,
		CacheWriteTokens: resp.Usage.CacheWriteTokens,
		CacheReadTokens:  resp.Usage.CacheReadTokens,
		Timestamp:        time.Now(),
	})

	if err != nil {
		log.Printf("[ai] summary generation error: %v", err)
		return
	}

	if err := postgre.UpdateTicketSummary(ctx, ticketID, resp.Text); err != nil {
		log.Printf("[ai] update summary error: %v", err)
		return
	}

	log.Printf("[ai] summary generated for ticket %s latency=%dms", ticketID, latency)
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
		Provider:         providerName,
		Model:            core.GetString(providerName+".model", "unknown"),
		LatencyMs:        latency,
		InputTokens:      resp.Usage.InputTokens,
		OutputTokens:     resp.Usage.OutputTokens,
		CacheWriteTokens: resp.Usage.CacheWriteTokens,
		CacheReadTokens:  resp.Usage.CacheReadTokens,
		Timestamp:        time.Now(),
	})

	if err != nil {
		return nil, err
	}

	var analysis structs.AIAnalysis
	text := strings.TrimSpace(resp.Text)
	if strings.HasPrefix(text, "```") {
		lines := strings.Split(text, "\n")
		if len(lines) > 2 {
			text = strings.Join(lines[1:len(lines)-1], "\n")
		}
	}
	if err := json.Unmarshal([]byte(text), &analysis); err != nil {
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
