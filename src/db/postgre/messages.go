package postgre

import (
	"context"

	"supportflow/core/structs"
)

func CreateMessage(ctx context.Context, m *structs.Message) error {
	return Pool.QueryRow(ctx,
		`INSERT INTO supportflow.messages (ticket_id, role, content)
		 VALUES ($1, $2, $3)
		 RETURNING id, created_at`,
		m.TicketID, m.Role, m.Content,
	).Scan(&m.ID, &m.CreatedAt)
}

func GetMessagesByTicket(ctx context.Context, ticketID string) ([]structs.Message, error) {
	rows, err := Pool.Query(ctx,
		`SELECT id, ticket_id, role, content, created_at
		 FROM supportflow.messages WHERE ticket_id = $1 ORDER BY created_at ASC`, ticketID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []structs.Message
	for rows.Next() {
		var m structs.Message
		if err := rows.Scan(&m.ID, &m.TicketID, &m.Role, &m.Content, &m.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, nil
}
