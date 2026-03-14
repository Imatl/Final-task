package postgre

import (
	"context"
	"fmt"
)

func CreateInviteToken(ctx context.Context, createdBy string) (string, error) {
	if Pool == nil {
		return "", fmt.Errorf("database pool is not initialized")
	}
	var token string
	err := Pool.QueryRow(ctx,
		`INSERT INTO supportflow.invite_tokens (created_by)
		 VALUES ($1)
		 RETURNING token`, createdBy,
	).Scan(&token)
	if err != nil {
		return "", err
	}
	return token, nil
}

func ValidateInviteToken(ctx context.Context, token string) error {
	if Pool == nil {
		return fmt.Errorf("database pool is not initialized")
	}
	var id string
	err := Pool.QueryRow(ctx,
		`SELECT id FROM supportflow.invite_tokens
		 WHERE token = $1 AND used = false AND expires_at > now()`, token,
	).Scan(&id)
	return err
}

func ConsumeInviteToken(ctx context.Context, token, usedBy string) error {
	if Pool == nil {
		return fmt.Errorf("database pool is not initialized")
	}
	tag, err := Pool.Exec(ctx,
		`UPDATE supportflow.invite_tokens
		 SET used = true, used_by = $2
		 WHERE token = $1 AND used = false AND expires_at > now()`, token, usedBy)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("invite token is invalid, expired or already used")
	}
	return nil
}
