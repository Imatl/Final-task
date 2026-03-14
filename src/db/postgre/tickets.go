package postgre

import (
	"context"
	"fmt"
	"strings"
	"time"

	"supportflow/core/structs"
)

func CreateTicket(ctx context.Context, t *structs.Ticket) error {
	return Pool.QueryRow(ctx,
		`INSERT INTO supportflow.tickets (customer_id, subject, channel, status, priority, category)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id, created_at, updated_at`,
		t.CustomerID, t.Subject, t.Channel, t.Status, t.Priority, t.Category,
	).Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
}

func GetTicket(ctx context.Context, id string) (*structs.Ticket, error) {
	if Pool == nil {
		return nil, fmt.Errorf("database pool is not initialized")
	}
	t := &structs.Ticket{}
	err := Pool.QueryRow(ctx,
		`SELECT id, customer_id, subject, channel, status, priority, category, agent_id, ai_summary, created_at, updated_at, closed_at
		 FROM supportflow.tickets WHERE id = $1`, id,
	).Scan(&t.ID, &t.CustomerID, &t.Subject, &t.Channel, &t.Status, &t.Priority, &t.Category, &t.AgentID, &t.AISummary, &t.CreatedAt, &t.UpdatedAt, &t.ClosedAt)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func UpdateTicketStatus(ctx context.Context, id, status string) error {
	if Pool == nil {
		return fmt.Errorf("database pool is not initialized")
	}
	var closedAt *time.Time
	if status == "resolved" || status == "closed" {
		now := time.Now()
		closedAt = &now
	}
	_, err := Pool.Exec(ctx,
		`UPDATE supportflow.tickets SET status = $1, updated_at = now(), closed_at = COALESCE($3, closed_at) WHERE id = $2`,
		status, id, closedAt,
	)
	return err
}

func UpdateTicketSummary(ctx context.Context, id, summary string) error {
	_, err := Pool.Exec(ctx,
		`UPDATE supportflow.tickets SET ai_summary = $1, updated_at = now() WHERE id = $2`,
		summary, id,
	)
	return err
}

func AssignTicket(ctx context.Context, ticketID, agentID string) error {
	_, err := Pool.Exec(ctx,
		`UPDATE supportflow.tickets SET agent_id = $1, status = 'in_progress', updated_at = now() WHERE id = $2`,
		agentID, ticketID,
	)
	return err
}

func ListTickets(ctx context.Context, f structs.TicketFilter) ([]structs.Ticket, int, error) {
	where := []string{"1=1"}
	args := []any{}
	idx := 1

	if f.Status != "" {
		where = append(where, fmt.Sprintf("status = $%d", idx))
		args = append(args, f.Status)
		idx++
	}
	if f.Priority != "" {
		where = append(where, fmt.Sprintf("priority = $%d", idx))
		args = append(args, f.Priority)
		idx++
	}
	if f.AgentID != "" {
		where = append(where, fmt.Sprintf("agent_id = $%d", idx))
		args = append(args, f.AgentID)
		idx++
	}
	if f.Category != "" {
		where = append(where, fmt.Sprintf("category = $%d", idx))
		args = append(args, f.Category)
		idx++
	}

	whereClause := strings.Join(where, " AND ")

	var total int
	err := Pool.QueryRow(ctx, "SELECT COUNT(*) FROM supportflow.tickets WHERE "+whereClause, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	limit := f.Limit
	if limit <= 0 {
		limit = 50
	}
	offset := f.Offset

	query := fmt.Sprintf(
		`SELECT id, customer_id, subject, channel, status, priority, category, agent_id, ai_summary, created_at, updated_at, closed_at
		 FROM supportflow.tickets WHERE %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d`,
		whereClause, idx, idx+1,
	)
	args = append(args, limit, offset)

	rows, err := Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var tickets []structs.Ticket
	for rows.Next() {
		var t structs.Ticket
		if err := rows.Scan(&t.ID, &t.CustomerID, &t.Subject, &t.Channel, &t.Status, &t.Priority, &t.Category, &t.AgentID, &t.AISummary, &t.CreatedAt, &t.UpdatedAt, &t.ClosedAt); err != nil {
			return nil, 0, err
		}
		tickets = append(tickets, t)
	}
	return tickets, total, nil
}
