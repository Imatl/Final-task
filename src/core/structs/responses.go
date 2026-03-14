package structs

type ChatResponse struct {
	TicketID  string      `json:"ticket_id"`
	Message   string      `json:"message"`
	Analysis  *AIAnalysis `json:"analysis,omitempty"`
	Actions   []Action    `json:"actions,omitempty"`
	AutoFixed bool        `json:"auto_fixed"`
}

type TicketDetail struct {
	Ticket   Ticket     `json:"ticket"`
	Customer Customer   `json:"customer"`
	Messages []Message  `json:"messages"`
	Actions  []Action   `json:"actions"`
	Analysis *AIAnalysis `json:"analysis,omitempty"`
}

type TicketListResponse struct {
	Tickets []Ticket `json:"tickets"`
	Total   int      `json:"total"`
}

type AnalyticsOverview struct {
	TotalTickets    int            `json:"total_tickets"`
	OpenTickets     int            `json:"open_tickets"`
	AvgResolveTime  float64        `json:"avg_resolve_time_minutes"`
	AutoResolveRate float64        `json:"auto_resolve_rate"`
	ByCategory      map[string]int `json:"by_category"`
	ByPriority      map[string]int `json:"by_priority"`
	BySentiment     map[string]int `json:"by_sentiment"`
}

type AgentPerformance struct {
	AgentID         string  `json:"agent_id"`
	AgentName       string  `json:"agent_name"`
	TicketsResolved int     `json:"tickets_resolved"`
	AvgResolveTime  float64 `json:"avg_resolve_time_minutes"`
	QualityScore    float64 `json:"quality_score"`
}
