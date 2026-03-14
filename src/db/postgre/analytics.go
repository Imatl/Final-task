package postgre

import (
	"context"

	"supportflow/core/structs"
)

func GetAnalyticsOverview(ctx context.Context) (*structs.AnalyticsOverview, error) {
	a := &structs.AnalyticsOverview{
		ByCategory:  make(map[string]int),
		ByPriority:  make(map[string]int),
		BySentiment: make(map[string]int),
	}

	Pool.QueryRow(ctx, `SELECT COUNT(*) FROM supportflow.tickets`).Scan(&a.TotalTickets)
	Pool.QueryRow(ctx, `SELECT COUNT(*) FROM supportflow.tickets WHERE status IN ('open','in_progress','waiting')`).Scan(&a.OpenTickets)
	Pool.QueryRow(ctx,
		`SELECT COALESCE(AVG(EXTRACT(EPOCH FROM (closed_at - created_at)) / 60), 0)
		 FROM supportflow.tickets WHERE closed_at IS NOT NULL`,
	).Scan(&a.AvgResolveTime)

	var autoResolved, totalResolved int
	Pool.QueryRow(ctx, `SELECT COUNT(*) FROM supportflow.tickets WHERE status IN ('resolved','closed')`).Scan(&totalResolved)
	Pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM supportflow.tickets t
		 WHERE t.status IN ('resolved','closed') AND t.agent_id IS NULL`,
	).Scan(&autoResolved)
	if totalResolved > 0 {
		a.AutoResolveRate = float64(autoResolved) / float64(totalResolved)
	}

	rows, _ := Pool.Query(ctx, `SELECT category, COUNT(*) FROM supportflow.tickets GROUP BY category`)
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var cat string
			var cnt int
			rows.Scan(&cat, &cnt)
			a.ByCategory[cat] = cnt
		}
	}

	rows2, _ := Pool.Query(ctx, `SELECT priority, COUNT(*) FROM supportflow.tickets GROUP BY priority`)
	if rows2 != nil {
		defer rows2.Close()
		for rows2.Next() {
			var pri string
			var cnt int
			rows2.Scan(&pri, &cnt)
			a.ByPriority[pri] = cnt
		}
	}

	rows3, _ := Pool.Query(ctx, `SELECT sentiment, COUNT(*) FROM supportflow.ai_analyses GROUP BY sentiment`)
	if rows3 != nil {
		defer rows3.Close()
		for rows3.Next() {
			var sent string
			var cnt int
			rows3.Scan(&sent, &cnt)
			a.BySentiment[sent] = cnt
		}
	}

	return a, nil
}

func GetAgentPerformance(ctx context.Context) ([]structs.AgentPerformance, error) {
	rows, err := Pool.Query(ctx,
		`SELECT a.id, a.name,
			COUNT(t.id) AS tickets_resolved,
			COALESCE(AVG(EXTRACT(EPOCH FROM (t.closed_at - t.created_at)) / 60), 0) AS avg_resolve_time
		 FROM supportflow.agents a
		 LEFT JOIN supportflow.tickets t ON t.agent_id = a.id AND t.status IN ('resolved','closed')
		 GROUP BY a.id, a.name
		 ORDER BY tickets_resolved DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var perfs []structs.AgentPerformance
	for rows.Next() {
		var p structs.AgentPerformance
		if err := rows.Scan(&p.AgentID, &p.AgentName, &p.TicketsResolved, &p.AvgResolveTime); err != nil {
			return nil, err
		}
		perfs = append(perfs, p)
	}
	return perfs, nil
}
