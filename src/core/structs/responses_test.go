package structs

import (
	"encoding/json"
	"testing"
)

func TestChatResponseJSON(t *testing.T) {
	resp := ChatResponse{
		TicketID:  "t-1",
		Message:   "Refund processed",
		AutoFixed: true,
		Actions: []Action{
			{ID: "a-1", Type: "refund", Status: "executed"},
		},
	}

	data, _ := json.Marshal(resp)
	var decoded ChatResponse
	json.Unmarshal(data, &decoded)

	if decoded.TicketID != "t-1" {
		t.Errorf("TicketID = %q, want t-1", decoded.TicketID)
	}
	if !decoded.AutoFixed {
		t.Error("AutoFixed should be true")
	}
	if len(decoded.Actions) != 1 {
		t.Errorf("Actions len = %d, want 1", len(decoded.Actions))
	}
}

func TestChatResponseOmitsNilAnalysis(t *testing.T) {
	resp := ChatResponse{TicketID: "t-2", Message: "ok"}

	data, _ := json.Marshal(resp)
	var raw map[string]any
	json.Unmarshal(data, &raw)

	if _, ok := raw["analysis"]; ok {
		t.Error("analysis should be omitted when nil")
	}
}

func TestTicketListResponseJSON(t *testing.T) {
	resp := TicketListResponse{
		Tickets: []Ticket{{ID: "t-1", Status: "open"}, {ID: "t-2", Status: "closed"}},
		Total:   42,
	}

	data, _ := json.Marshal(resp)
	var decoded TicketListResponse
	json.Unmarshal(data, &decoded)

	if decoded.Total != 42 {
		t.Errorf("Total = %d, want 42", decoded.Total)
	}
	if len(decoded.Tickets) != 2 {
		t.Errorf("Tickets len = %d, want 2", len(decoded.Tickets))
	}
}

func TestAnalyticsOverviewJSON(t *testing.T) {
	overview := AnalyticsOverview{
		TotalTickets:    100,
		OpenTickets:     25,
		AvgResolveTime:  45.5,
		AutoResolveRate: 0.72,
		ByCategory:      map[string]int{"billing": 30, "technical": 40},
		ByPriority:      map[string]int{"high": 20, "low": 50},
		BySentiment:     map[string]int{"angry": 10, "neutral": 60},
	}

	data, _ := json.Marshal(overview)
	var decoded AnalyticsOverview
	json.Unmarshal(data, &decoded)

	if decoded.TotalTickets != 100 {
		t.Errorf("TotalTickets = %d, want 100", decoded.TotalTickets)
	}
	if decoded.ByCategory["billing"] != 30 {
		t.Errorf("ByCategory[billing] = %d, want 30", decoded.ByCategory["billing"])
	}
	if decoded.AutoResolveRate != 0.72 {
		t.Errorf("AutoResolveRate = %f, want 0.72", decoded.AutoResolveRate)
	}
}

func TestAgentPerformanceJSON(t *testing.T) {
	perf := AgentPerformance{
		AgentID:         "ag-1",
		AgentName:       "Alice",
		TicketsResolved: 50,
		AvgResolveTime:  30.2,
		QualityScore:    4.8,
	}

	data, _ := json.Marshal(perf)
	var decoded AgentPerformance
	json.Unmarshal(data, &decoded)

	if decoded.AgentName != "Alice" {
		t.Errorf("AgentName = %q, want Alice", decoded.AgentName)
	}
	if decoded.QualityScore != 4.8 {
		t.Errorf("QualityScore = %f, want 4.8", decoded.QualityScore)
	}
}

func TestTicketDetailJSON(t *testing.T) {
	detail := TicketDetail{
		Ticket:   Ticket{ID: "t-1", Status: "open"},
		Customer: Customer{ID: "c-1", Name: "Bob"},
		Messages: []Message{{ID: "m-1", Content: "help"}},
		Actions:  []Action{{ID: "a-1", Type: "refund"}},
	}

	data, _ := json.Marshal(detail)
	var decoded TicketDetail
	json.Unmarshal(data, &decoded)

	if decoded.Ticket.ID != "t-1" {
		t.Errorf("Ticket.ID = %q, want t-1", decoded.Ticket.ID)
	}
	if decoded.Customer.Name != "Bob" {
		t.Errorf("Customer.Name = %q, want Bob", decoded.Customer.Name)
	}
	if decoded.Analysis != nil {
		t.Error("Analysis should be nil")
	}
}
