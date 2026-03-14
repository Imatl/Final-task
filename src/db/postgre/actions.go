package postgre

import (
	"context"
	"time"

	"supportflow/core/structs"
)

func CreateAction(ctx context.Context, a *structs.Action) error {
	return Pool.QueryRow(ctx,
		`INSERT INTO supportflow.actions (ticket_id, type, params, status, confidence)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, created_at`,
		a.TicketID, a.Type, a.Params, a.Status, a.Confidence,
	).Scan(&a.ID, &a.CreatedAt)
}

func GetAction(ctx context.Context, id string) (*structs.Action, error) {
	a := &structs.Action{}
	err := Pool.QueryRow(ctx,
		`SELECT id, ticket_id, type, params, status, result, confidence, created_at, executed_at
		 FROM supportflow.actions WHERE id = $1`, id,
	).Scan(&a.ID, &a.TicketID, &a.Type, &a.Params, &a.Status, &a.Result, &a.Confidence, &a.CreatedAt, &a.ExecutedAt)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func GetActionsByTicket(ctx context.Context, ticketID string) ([]structs.Action, error) {
	rows, err := Pool.Query(ctx,
		`SELECT id, ticket_id, type, params, status, result, confidence, created_at, executed_at
		 FROM supportflow.actions WHERE ticket_id = $1 ORDER BY created_at ASC`, ticketID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var actions []structs.Action
	for rows.Next() {
		var a structs.Action
		if err := rows.Scan(&a.ID, &a.TicketID, &a.Type, &a.Params, &a.Status, &a.Result, &a.Confidence, &a.CreatedAt, &a.ExecutedAt); err != nil {
			return nil, err
		}
		actions = append(actions, a)
	}
	return actions, nil
}

func UpdateActionStatus(ctx context.Context, id, status, result string) error {
	var executedAt *time.Time
	if status == "executed" || status == "approved" {
		now := time.Now()
		executedAt = &now
	}
	_, err := Pool.Exec(ctx,
		`UPDATE supportflow.actions SET status = $1, result = $2, executed_at = $3 WHERE id = $4`,
		status, result, executedAt, id,
	)
	return err
}

func GetPendingActions(ctx context.Context) ([]structs.Action, error) {
	rows, err := Pool.Query(ctx,
		`SELECT id, ticket_id, type, params, status, result, confidence, created_at, executed_at
		 FROM supportflow.actions WHERE status = 'pending' ORDER BY created_at ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var actions []structs.Action
	for rows.Next() {
		var a structs.Action
		if err := rows.Scan(&a.ID, &a.TicketID, &a.Type, &a.Params, &a.Status, &a.Result, &a.Confidence, &a.CreatedAt, &a.ExecutedAt); err != nil {
			return nil, err
		}
		actions = append(actions, a)
	}
	return actions, nil
}
