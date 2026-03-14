package structs

type ChatRequest struct {
	TicketID   string `json:"ticket_id,omitempty"`
	CustomerID string `json:"customer_id"`
	Channel    string `json:"channel,omitempty"`
	Message    string `json:"message"`
}

type ActionApproval struct {
	ActionID string `json:"action_id"`
	Approved bool   `json:"approved"`
	AgentID  string `json:"agent_id"`
}

type TicketFilter struct {
	Status   string `json:"status,omitempty"`
	Priority string `json:"priority,omitempty"`
	AgentID  string `json:"agent_id,omitempty"`
	Category string `json:"category,omitempty"`
	Limit    int    `json:"limit,omitempty"`
	Offset   int    `json:"offset,omitempty"`
}
