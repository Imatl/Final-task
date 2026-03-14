package postgre

import (
	"context"
	"fmt"

	"supportflow/core/structs"
)

func GetCustomer(ctx context.Context, id string) (*structs.Customer, error) {
	if Pool == nil {
		return nil, fmt.Errorf("database pool is not initialized")
	}
	c := &structs.Customer{}
	err := Pool.QueryRow(ctx,
		`SELECT id, name, email, phone, plan, company, created_at FROM supportflow.customers WHERE id = $1`, id,
	).Scan(&c.ID, &c.Name, &c.Email, &c.Phone, &c.Plan, &c.Company, &c.CreatedAt)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func FindOrCreateCustomer(ctx context.Context, email, name, plan, phone string) (*structs.Customer, error) {
	if Pool == nil {
		return nil, fmt.Errorf("database pool is not initialized")
	}
	c := &structs.Customer{}
	err := Pool.QueryRow(ctx,
		`SELECT id, name, email, phone, plan, company, created_at FROM supportflow.customers WHERE email = $1`, email,
	).Scan(&c.ID, &c.Name, &c.Email, &c.Phone, &c.Plan, &c.Company, &c.CreatedAt)
	if err == nil {
		return c, nil
	}

	var phoneVal *string
	if phone != "" {
		phoneVal = &phone
	}

	err = Pool.QueryRow(ctx,
		`INSERT INTO supportflow.customers (name, email, phone, plan) VALUES ($1, $2, $3, $4)
		 RETURNING id, name, email, phone, plan, company, created_at`,
		name, email, phoneVal, plan,
	).Scan(&c.ID, &c.Name, &c.Email, &c.Phone, &c.Plan, &c.Company, &c.CreatedAt)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func ListCustomers(ctx context.Context) ([]structs.Customer, error) {
	rows, err := Pool.Query(ctx,
		`SELECT id, name, email, phone, plan, company, created_at FROM supportflow.customers ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []structs.Customer
	for rows.Next() {
		var c structs.Customer
		if err := rows.Scan(&c.ID, &c.Name, &c.Email, &c.Phone, &c.Plan, &c.Company, &c.CreatedAt); err != nil {
			return nil, err
		}
		customers = append(customers, c)
	}
	return customers, nil
}
