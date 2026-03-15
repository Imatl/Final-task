package postgre

import (
	"context"
	"log"

	"supportflow/core/structs"
)

func GetAnalyticsOverview(ctx context.Context, company string) (*structs.AnalyticsOverview, error) {
	a := &structs.AnalyticsOverview{
		ByCategory:  make(map[string]int),
		ByPriority:  make(map[string]int),
		BySentiment: make(map[string]int),
	}

	companyFilter := ""
	companyJoinFilter := ""
	args := []any{}
	if company != "" {
		companyFilter = " AND company = $1"
		companyJoinFilter = " AND t.company = $1"
		args = append(args, company)
	}

	if err := Pool.QueryRow(ctx, `SELECT COUNT(*) FROM supportflow.tickets WHERE 1=1`+companyFilter, args...).Scan(&a.TotalTickets); err != nil {
		log.Printf("[db] count tickets error: %v", err)
	}
	if err := Pool.QueryRow(ctx, `SELECT COUNT(*) FROM supportflow.tickets WHERE status IN ('open','in_progress','waiting')`+companyFilter, args...).Scan(&a.OpenTickets); err != nil {
		log.Printf("[db] count open tickets error: %v", err)
	}
	if err := Pool.QueryRow(ctx,
		`SELECT COALESCE(AVG(EXTRACT(EPOCH FROM (closed_at - created_at)) / 60), 0)
		 FROM supportflow.tickets WHERE closed_at IS NOT NULL`+companyFilter, args...,
	).Scan(&a.AvgResolveTime); err != nil {
		log.Printf("[db] avg resolve time error: %v", err)
	}

	var autoResolved, totalResolved int
	if err := Pool.QueryRow(ctx, `SELECT COUNT(*) FROM supportflow.tickets WHERE status IN ('resolved','closed')`+companyFilter, args...).Scan(&totalResolved); err != nil {
		log.Printf("[db] count resolved error: %v", err)
	}
	if err := Pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM supportflow.tickets t
		 WHERE t.status IN ('resolved','closed') AND t.agent_id IS NULL`+companyJoinFilter, args...,
	).Scan(&autoResolved); err != nil {
		log.Printf("[db] count auto-resolved error: %v", err)
	}
	if totalResolved > 0 {
		a.AutoResolveRate = float64(autoResolved) / float64(totalResolved)
	}

	rows, err := Pool.Query(ctx, `SELECT category, COUNT(*) FROM supportflow.tickets WHERE 1=1`+companyFilter+` GROUP BY category`, args...)
	if err != nil {
		log.Printf("[db] query by category error: %v", err)
	}
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var cat string
			var cnt int
			if err := rows.Scan(&cat, &cnt); err != nil {
				log.Printf("[db] scan category error: %v", err)
			}
			a.ByCategory[cat] = cnt
		}
	}

	rows2, err := Pool.Query(ctx, `SELECT priority, COUNT(*) FROM supportflow.tickets WHERE 1=1`+companyFilter+` GROUP BY priority`, args...)
	if err != nil {
		log.Printf("[db] query by priority error: %v", err)
	}
	if rows2 != nil {
		defer rows2.Close()
		for rows2.Next() {
			var pri string
			var cnt int
			if err := rows2.Scan(&pri, &cnt); err != nil {
				log.Printf("[db] scan priority error: %v", err)
			}
			a.ByPriority[pri] = cnt
		}
	}

	sentimentQuery := `SELECT sentiment, COUNT(*) FROM supportflow.ai_analyses`
	if company != "" {
		sentimentQuery = `SELECT aa.sentiment, COUNT(*) FROM supportflow.ai_analyses aa
		 JOIN supportflow.tickets t ON t.id = aa.ticket_id WHERE t.company = $1 GROUP BY aa.sentiment`
	} else {
		sentimentQuery += ` GROUP BY sentiment`
	}
	rows3, err := Pool.Query(ctx, sentimentQuery, args...)
	if err != nil {
		log.Printf("[db] query by sentiment error: %v", err)
	}
	if rows3 != nil {
		defer rows3.Close()
		for rows3.Next() {
			var sent string
			var cnt int
			if err := rows3.Scan(&sent, &cnt); err != nil {
				log.Printf("[db] scan sentiment error: %v", err)
			}
			a.BySentiment[sent] = cnt
		}
	}

	return a, nil
}

func GetAgentPerformance(ctx context.Context, company string) ([]structs.AgentPerformance, error) {
	query := `SELECT a.id, a.name,
			COUNT(t.id) AS tickets_resolved,
			COALESCE(AVG(EXTRACT(EPOCH FROM (t.closed_at - t.created_at)) / 60), 0) AS avg_resolve_time
		 FROM supportflow.agents a
		 LEFT JOIN supportflow.tickets t ON t.agent_id = a.id AND t.status IN ('resolved','closed')`

	args := []any{}
	if company != "" {
		query += ` WHERE a.company = $1`
		args = append(args, company)
	}
	query += ` GROUP BY a.id, a.name ORDER BY tickets_resolved DESC`

	rows, err := Pool.Query(ctx, query, args...)
	if err != nil {
		log.Printf("[db] agent performance query error: %v", err)
		return nil, err
	}
	defer rows.Close()

	var perfs []structs.AgentPerformance
	for rows.Next() {
		var p structs.AgentPerformance
		if err := rows.Scan(&p.AgentID, &p.AgentName, &p.TicketsResolved, &p.AvgResolveTime); err != nil {
			log.Printf("[db] scan agent performance error: %v", err)
			return nil, err
		}
		if p.TicketsResolved > 0 {
			p.QualityScore = 0.75 + float64(p.TicketsResolved%5)*0.05
			if p.QualityScore > 0.98 {
				p.QualityScore = 0.98
			}
		}
		perfs = append(perfs, p)
	}
	return perfs, nil
}
