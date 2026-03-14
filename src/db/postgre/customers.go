package postgre

import (
	"context"

	"supportflow/core/structs"
)

func GetCustomer(ctx context.Context, id string) (*structs.Customer, error) {
	c := &structs.Customer{}
	err := Pool.QueryRow(ctx,
		`SELECT id, name, email, plan, created_at FROM supportflow.customers WHERE id = $1`, id,
	).Scan(&c.ID, &c.Name, &c.Email, &c.Plan, &c.CreatedAt)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func ListCustomers(ctx context.Context) ([]structs.Customer, error) {
	rows, err := Pool.Query(ctx,
		`SELECT id, name, email, plan, created_at FROM supportflow.customers ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []structs.Customer
	for rows.Next() {
		var c structs.Customer
		if err := rows.Scan(&c.ID, &c.Name, &c.Email, &c.Plan, &c.CreatedAt); err != nil {
			return nil, err
		}
		customers = append(customers, c)
	}
	return customers, nil
}
