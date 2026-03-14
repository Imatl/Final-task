package structs

import "time"

type Customer struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     *string   `json:"phone,omitempty"`
	Plan      string    `json:"plan"`
	Company   *string   `json:"company,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type Ticket struct {
	ID         string     `json:"id"`
	CustomerID string     `json:"customer_id"`
	Subject    string     `json:"subject"`
	Channel    string     `json:"channel"`
	Status     string     `json:"status"`
	Priority   string     `json:"priority"`
	Category   string     `json:"category"`
	AgentID    *string    `json:"agent_id,omitempty"`
	AISummary  *string    `json:"ai_summary,omitempty"`
	Company    *string    `json:"company,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	ClosedAt   *time.Time `json:"closed_at,omitempty"`
}

type Message struct {
	ID        string    `json:"id"`
	TicketID  string    `json:"ticket_id"`
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type Action struct {
	ID         string     `json:"id"`
	TicketID   string     `json:"ticket_id"`
	Type       string     `json:"type"`
	Params     string     `json:"params"`
	Status     string     `json:"status"`
	Result     *string    `json:"result,omitempty"`
	Confidence float64    `json:"confidence"`
	CreatedAt  time.Time  `json:"created_at"`
	ExecutedAt *time.Time `json:"executed_at,omitempty"`
}

type Agent struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	IsOnline  bool      `json:"is_online"`
	Company   *string   `json:"company,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type AIAnalysis struct {
	TicketID       string   `json:"ticket_id"`
	Intent         string   `json:"intent"`
	Sentiment      string   `json:"sentiment"`
	Urgency        string   `json:"urgency"`
	SuggestedTools []string `json:"suggested_tools"`
	Reasoning      string   `json:"reasoning"`
	Confidence     float64  `json:"confidence"`
}

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	GoogleSub *string   `json:"google_sub,omitempty"`
	Password  *string   `json:"-"`
	Level     int       `json:"level"`
	Role      string    `json:"role"`
	Company   *string   `json:"company,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type KBEntry struct {
	ID        string    `json:"id"`
	Company   string    `json:"company"`
	Question  string    `json:"question"`
	Answer    string    `json:"answer"`
	CreatedAt time.Time `json:"created_at"`
}
