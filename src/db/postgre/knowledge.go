package postgre

import (
	"context"
	"fmt"

	"supportflow/core/structs"
)

func ListKBEntries(ctx context.Context, company string) ([]structs.KBEntry, error) {
	if Pool == nil {
		return nil, fmt.Errorf("database pool is not initialized")
	}
	rows, err := Pool.Query(ctx,
		`SELECT id, company, question, answer, created_at
		 FROM supportflow.knowledge_base WHERE company = $1
		 ORDER BY created_at DESC`, company)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []structs.KBEntry
	for rows.Next() {
		var e structs.KBEntry
		if err := rows.Scan(&e.ID, &e.Company, &e.Question, &e.Answer, &e.CreatedAt); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, nil
}

func CreateKBEntry(ctx context.Context, company, question, answer string) (*structs.KBEntry, error) {
	if Pool == nil {
		return nil, fmt.Errorf("database pool is not initialized")
	}
	e := &structs.KBEntry{}
	err := Pool.QueryRow(ctx,
		`INSERT INTO supportflow.knowledge_base (company, question, answer)
		 VALUES ($1, $2, $3)
		 RETURNING id, company, question, answer, created_at`,
		company, question, answer,
	).Scan(&e.ID, &e.Company, &e.Question, &e.Answer, &e.CreatedAt)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func DeleteKBEntry(ctx context.Context, id string) error {
	if Pool == nil {
		return fmt.Errorf("database pool is not initialized")
	}
	_, err := Pool.Exec(ctx,
		`DELETE FROM supportflow.knowledge_base WHERE id = $1`, id)
	return err
}
