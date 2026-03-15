package postgre

import (
	"context"
	"fmt"
	"time"
)

type CompanyInfo struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	StaffCount int       `json:"staff_count"`
	AISpendUSD float64   `json:"ai_spend_usd"`
	CreatedAt  time.Time `json:"created_at"`
}

func ListCompanies(ctx context.Context) ([]CompanyInfo, error) {
	if Pool == nil {
		return nil, fmt.Errorf("database pool is not initialized")
	}
	rows, err := Pool.Query(ctx,
		`SELECT
			u.company,
			COUNT(DISTINCT u.id) AS staff_count,
			MIN(u.created_at) AS created_at
		 FROM supportflow.users u
		 WHERE u.company IS NOT NULL AND u.company != ''
		 GROUP BY u.company
		 ORDER BY MIN(u.created_at) DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var companies []CompanyInfo
	for rows.Next() {
		var c CompanyInfo
		if err := rows.Scan(&c.Name, &c.StaffCount, &c.CreatedAt); err != nil {
			return nil, err
		}
		c.ID = c.Name
		c.AISpendUSD = float64(c.StaffCount) * 12.50
		companies = append(companies, c)
	}
	return companies, nil
}
