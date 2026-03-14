package postgre

import (
	"context"
	"fmt"

	"supportflow/core/structs"
)

func FindUserByEmail(ctx context.Context, email string) (*structs.User, error) {
	if Pool == nil {
		return nil, fmt.Errorf("database pool is not initialized")
	}
	u := &structs.User{}
	err := Pool.QueryRow(ctx,
		`SELECT id, email, name, google_sub, password, level, role, company, created_at
		 FROM supportflow.users WHERE email = $1`, email,
	).Scan(&u.ID, &u.Email, &u.Name, &u.GoogleSub, &u.Password, &u.Level, &u.Role, &u.Company, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func FindUserByGoogleSub(ctx context.Context, sub string) (*structs.User, error) {
	if Pool == nil {
		return nil, fmt.Errorf("database pool is not initialized")
	}
	u := &structs.User{}
	err := Pool.QueryRow(ctx,
		`SELECT id, email, name, google_sub, password, level, role, company, created_at
		 FROM supportflow.users WHERE google_sub = $1`, sub,
	).Scan(&u.ID, &u.Email, &u.Name, &u.GoogleSub, &u.Password, &u.Level, &u.Role, &u.Company, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func CreateUser(ctx context.Context, email, name, googleSub string) (*structs.User, error) {
	if Pool == nil {
		return nil, fmt.Errorf("database pool is not initialized")
	}
	u := &structs.User{}
	var sub *string
	if googleSub != "" {
		sub = &googleSub
	}
	err := Pool.QueryRow(ctx,
		`INSERT INTO supportflow.users (email, name, google_sub)
		 VALUES ($1, $2, $3)
		 RETURNING id, email, name, google_sub, password, level, role, company, created_at`,
		email, name, sub,
	).Scan(&u.ID, &u.Email, &u.Name, &u.GoogleSub, &u.Password, &u.Level, &u.Role, &u.Company, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func CreateUserWithPassword(ctx context.Context, email, name, password, company string, level int, role string) (*structs.User, error) {
	if Pool == nil {
		return nil, fmt.Errorf("database pool is not initialized")
	}
	u := &structs.User{}
	err := Pool.QueryRow(ctx,
		`INSERT INTO supportflow.users (email, name, password, company, level, role)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id, email, name, google_sub, password, level, role, company, created_at`,
		email, name, password, company, level, role,
	).Scan(&u.ID, &u.Email, &u.Name, &u.GoogleSub, &u.Password, &u.Level, &u.Role, &u.Company, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func UpdateUserGoogleSub(ctx context.Context, id, sub string) error {
	if Pool == nil {
		return fmt.Errorf("database pool is not initialized")
	}
	_, err := Pool.Exec(ctx,
		`UPDATE supportflow.users SET google_sub = $1 WHERE id = $2`, sub, id)
	return err
}
