package ai

import (
	"context"
	"testing"
)

func TestProcessMessageNoProvider(t *testing.T) {
	resetProviders()

	_, err := ProcessMessage(context.Background(), "ticket-1", "hello", "")
	if err == nil {
		t.Error("ProcessMessage should error when no provider configured")
	}
	if err.Error() != "no AI provider configured" {
		t.Errorf("error = %q, want 'no AI provider configured'", err.Error())
	}
}

func TestAnalyzeTicketNoProvider(t *testing.T) {
	resetProviders()

	_, err := AnalyzeTicket(context.Background(), "ticket-1", "hello")
	if err == nil {
		t.Error("AnalyzeTicket should error when no provider configured")
	}
}

func TestAnalyzeTicketValidJSON(t *testing.T) {
	resetProviders()

	mock := &mockProvider{
		name: "test",
		response: &LLMResponse{
			Text:  `{"intent":"billing_dispute","sentiment":"angry","urgency":"high","suggested_tools":["refund"],"reasoning":"double charge","confidence":0.95}`,
			Usage: TokenUsage{InputTokens: 30, OutputTokens: 40},
		},
	}
	RegisterProvider("test", mock)

	analysis, err := AnalyzeTicket(context.Background(), "ticket-1", "I was charged twice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if analysis.TicketID != "ticket-1" {
		t.Errorf("TicketID = %q, want ticket-1", analysis.TicketID)
	}
	if analysis.Intent != "billing_dispute" {
		t.Errorf("Intent = %q, want billing_dispute", analysis.Intent)
	}
	if analysis.Sentiment != "angry" {
		t.Errorf("Sentiment = %q, want angry", analysis.Sentiment)
	}
	if analysis.Confidence != 0.95 {
		t.Errorf("Confidence = %f, want 0.95", analysis.Confidence)
	}
	if len(analysis.SuggestedTools) != 1 || analysis.SuggestedTools[0] != "refund" {
		t.Errorf("SuggestedTools = %v, want [refund]", analysis.SuggestedTools)
	}
}

func TestAnalyzeTicketInvalidJSON(t *testing.T) {
	resetProviders()

	mock := &mockProvider{
		name: "test",
		response: &LLMResponse{
			Text:  "This is not JSON at all",
			Usage: TokenUsage{InputTokens: 10, OutputTokens: 15},
		},
	}
	RegisterProvider("test", mock)

	analysis, err := AnalyzeTicket(context.Background(), "ticket-1", "help")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if analysis.Intent != "general_inquiry" {
		t.Errorf("fallback Intent = %q, want general_inquiry", analysis.Intent)
	}
	if analysis.Sentiment != "neutral" {
		t.Errorf("fallback Sentiment = %q, want neutral", analysis.Sentiment)
	}
	if analysis.Confidence != 0.5 {
		t.Errorf("fallback Confidence = %f, want 0.5", analysis.Confidence)
	}
	if analysis.TicketID != "ticket-1" {
		t.Errorf("fallback TicketID = %q, want ticket-1", analysis.TicketID)
	}
}

func TestAnalyzeTicketPartialJSON(t *testing.T) {
	resetProviders()

	mock := &mockProvider{
		name: "test",
		response: &LLMResponse{
			Text:  `{"intent":"cancellation","sentiment":"negative","urgency":"medium","suggested_tools":[],"reasoning":"wants to leave","confidence":0.8}`,
			Usage: TokenUsage{InputTokens: 20, OutputTokens: 30},
		},
	}
	RegisterProvider("test", mock)

	analysis, err := AnalyzeTicket(context.Background(), "ticket-2", "I want to cancel")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if analysis.Intent != "cancellation" {
		t.Errorf("Intent = %q, want cancellation", analysis.Intent)
	}
	if analysis.TicketID != "ticket-2" {
		t.Errorf("TicketID = %q, want ticket-2 (should override)", analysis.TicketID)
	}
}

func TestAnalyzeTicketLLMError(t *testing.T) {
	resetProviders()

	mock := &mockProvider{
		name:     "test",
		response: &LLMResponse{Usage: TokenUsage{}},
		err:      nil,
	}
	RegisterProvider("test", mock)

	analysis, err := AnalyzeTicket(context.Background(), "ticket-1", "test")
	if err != nil {
		t.Fatalf("should not error for empty text: %v", err)
	}
	if analysis.Intent != "general_inquiry" {
		t.Errorf("empty response should fallback, got Intent = %q", analysis.Intent)
	}
}
