package structs

import (
	"encoding/json"
	"testing"
)

func TestChatRequestJSON(t *testing.T) {
	input := `{"customer_id":"c-1","message":"help me","ticket_id":"t-1","channel":"web"}`
	var req ChatRequest
	if err := json.Unmarshal([]byte(input), &req); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if req.CustomerID != "c-1" {
		t.Errorf("CustomerID = %q, want c-1", req.CustomerID)
	}
	if req.Message != "help me" {
		t.Errorf("Message = %q, want 'help me'", req.Message)
	}
	if req.TicketID != "t-1" {
		t.Errorf("TicketID = %q, want t-1", req.TicketID)
	}
	if req.Channel != "web" {
		t.Errorf("Channel = %q, want web", req.Channel)
	}
}

func TestChatRequestOmitEmpty(t *testing.T) {
	req := ChatRequest{
		CustomerID: "c-1",
		Message:    "help",
	}

	data, _ := json.Marshal(req)
	var raw map[string]any
	json.Unmarshal(data, &raw)

	if _, ok := raw["ticket_id"]; ok {
		t.Error("ticket_id should be omitted when empty")
	}
	if _, ok := raw["channel"]; ok {
		t.Error("channel should be omitted when empty")
	}
}

func TestActionApprovalJSON(t *testing.T) {
	input := `{"action_id":"a-1","approved":true,"agent_id":"ag-1"}`
	var approval ActionApproval
	if err := json.Unmarshal([]byte(input), &approval); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if approval.ActionID != "a-1" {
		t.Errorf("ActionID = %q, want a-1", approval.ActionID)
	}
	if !approval.Approved {
		t.Error("Approved should be true")
	}
	if approval.AgentID != "ag-1" {
		t.Errorf("AgentID = %q, want ag-1", approval.AgentID)
	}
}

func TestActionApprovalRejected(t *testing.T) {
	input := `{"action_id":"a-2","approved":false,"agent_id":"ag-2"}`
	var approval ActionApproval
	json.Unmarshal([]byte(input), &approval)

	if approval.Approved {
		t.Error("Approved should be false")
	}
}

func TestTicketFilterJSON(t *testing.T) {
	input := `{"status":"open","priority":"high","agent_id":"ag-1","category":"billing","limit":10,"offset":20}`
	var filter TicketFilter
	if err := json.Unmarshal([]byte(input), &filter); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if filter.Status != "open" {
		t.Errorf("Status = %q, want open", filter.Status)
	}
	if filter.Limit != 10 {
		t.Errorf("Limit = %d, want 10", filter.Limit)
	}
	if filter.Offset != 20 {
		t.Errorf("Offset = %d, want 20", filter.Offset)
	}
}

func TestTicketFilterOmitEmpty(t *testing.T) {
	filter := TicketFilter{Limit: 5}

	data, _ := json.Marshal(filter)
	var raw map[string]any
	json.Unmarshal(data, &raw)

	if _, ok := raw["status"]; ok {
		t.Error("status should be omitted when empty")
	}
	if _, ok := raw["priority"]; ok {
		t.Error("priority should be omitted when empty")
	}
}
