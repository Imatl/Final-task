package structs

import (
	"encoding/json"
	"testing"
	"time"
)

func ptr(s string) *string { return &s }

func TestTicketJSONRoundTrip(t *testing.T) {
	now := time.Date(2026, 3, 14, 12, 0, 0, 0, time.UTC)
	agentID := "agent-1"
	summary := "test summary"
	ticket := Ticket{
		ID:         "t-1",
		CustomerID: "c-1",
		Subject:    "Help",
		Channel:    "web",
		Status:     "open",
		Priority:   "high",
		Category:   "billing",
		AgentID:    &agentID,
		AISummary:  &summary,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	data, err := json.Marshal(ticket)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var decoded Ticket
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if decoded.ID != ticket.ID {
		t.Errorf("ID = %q, want %q", decoded.ID, ticket.ID)
	}
	if decoded.AgentID == nil || *decoded.AgentID != agentID {
		t.Errorf("AgentID = %v, want %q", decoded.AgentID, agentID)
	}
	if decoded.ClosedAt != nil {
		t.Error("ClosedAt should be nil")
	}
}

func TestTicketJSONOmitsNilFields(t *testing.T) {
	ticket := Ticket{
		ID:        "t-2",
		Status:    "open",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	data, _ := json.Marshal(ticket)
	var raw map[string]any
	json.Unmarshal(data, &raw)

	if _, ok := raw["agent_id"]; ok {
		t.Error("agent_id should be omitted when nil")
	}
	if _, ok := raw["ai_summary"]; ok {
		t.Error("ai_summary should be omitted when nil")
	}
	if _, ok := raw["closed_at"]; ok {
		t.Error("closed_at should be omitted when nil")
	}
}

func TestCustomerJSON(t *testing.T) {
	c := Customer{
		ID:    "c-1",
		Name:  "Alice",
		Email: "alice@test.com",
		Plan:  "premium",
	}

	data, _ := json.Marshal(c)
	var decoded Customer
	json.Unmarshal(data, &decoded)

	if decoded.Name != "Alice" {
		t.Errorf("Name = %q, want Alice", decoded.Name)
	}
	if decoded.Email != "alice@test.com" {
		t.Errorf("Email = %q, want alice@test.com", decoded.Email)
	}
}

func TestActionJSONWithOptionalFields(t *testing.T) {
	now := time.Now()
	a := Action{
		ID:         "a-1",
		TicketID:   "t-1",
		Type:       "refund",
		Params:     `{"amount":9.99}`,
		Status:     "executed",
		Result:     ptr(`{"success":true}`),
		Confidence: 0.95,
		CreatedAt:  now,
		ExecutedAt: &now,
	}

	data, _ := json.Marshal(a)
	var decoded Action
	json.Unmarshal(data, &decoded)

	if decoded.Confidence != 0.95 {
		t.Errorf("Confidence = %f, want 0.95", decoded.Confidence)
	}
	if decoded.ExecutedAt == nil {
		t.Error("ExecutedAt should not be nil")
	}
}

func TestActionJSONOmitsNilExecutedAt(t *testing.T) {
	a := Action{ID: "a-2", Status: "pending"}

	data, _ := json.Marshal(a)
	var raw map[string]any
	json.Unmarshal(data, &raw)

	if _, ok := raw["executed_at"]; ok {
		t.Error("executed_at should be omitted when nil")
	}
}

func TestMessageJSON(t *testing.T) {
	m := Message{
		ID:       "m-1",
		TicketID: "t-1",
		Role:     "customer",
		Content:  "I need help",
	}

	data, _ := json.Marshal(m)
	var decoded Message
	json.Unmarshal(data, &decoded)

	if decoded.Content != "I need help" {
		t.Errorf("Content = %q, want 'I need help'", decoded.Content)
	}
	if decoded.Role != "customer" {
		t.Errorf("Role = %q, want 'customer'", decoded.Role)
	}
}

func TestAIAnalysisJSON(t *testing.T) {
	a := AIAnalysis{
		TicketID:       "t-1",
		Intent:         "billing_dispute",
		Sentiment:      "angry",
		Urgency:        "high",
		SuggestedTools: []string{"refund", "lookup_billing"},
		Reasoning:      "double charge detected",
		Confidence:     0.92,
	}

	data, _ := json.Marshal(a)
	var decoded AIAnalysis
	json.Unmarshal(data, &decoded)

	if len(decoded.SuggestedTools) != 2 {
		t.Errorf("SuggestedTools len = %d, want 2", len(decoded.SuggestedTools))
	}
	if decoded.Confidence != 0.92 {
		t.Errorf("Confidence = %f, want 0.92", decoded.Confidence)
	}
}

func TestAgentJSON(t *testing.T) {
	a := Agent{
		ID:       "ag-1",
		Name:     "Bob",
		Email:    "bob@test.com",
		Role:     "agent",
		IsOnline: true,
	}

	data, _ := json.Marshal(a)
	var decoded Agent
	json.Unmarshal(data, &decoded)

	if !decoded.IsOnline {
		t.Error("IsOnline should be true")
	}
}
