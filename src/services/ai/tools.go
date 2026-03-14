package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"supportflow/core/constants"
	"supportflow/db/postgre"
)

type ToolResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func ExecuteTool(ctx context.Context, ticketID, toolName string, paramsJSON string) ToolResult {
	log.Printf("Executing tool: %s for ticket %s with params: %s", toolName, ticketID, paramsJSON)

	switch toolName {
	case constants.ActionRefund:
		return executeRefund(ctx, ticketID, paramsJSON)
	case constants.ActionChangePlan:
		return executeChangePlan(ctx, ticketID, paramsJSON)
	case constants.ActionResetPassword:
		return executeResetPassword(ctx, ticketID, paramsJSON)
	case constants.ActionEscalate:
		return executeEscalate(ctx, ticketID, paramsJSON)
	case constants.ActionSendEmail:
		return executeSendEmail(ctx, ticketID, paramsJSON)
	case constants.ActionCancelSub:
		return executeCancelSub(ctx, ticketID, paramsJSON)
	case "lookup_billing":
		return executeLookupBilling(ctx, ticketID, paramsJSON)
	case "lookup_customer":
		return executeLookupCustomer(ctx, ticketID, paramsJSON)
	default:
		return ToolResult{Success: false, Message: fmt.Sprintf("unknown tool: %s", toolName)}
	}
}

func executeRefund(ctx context.Context, ticketID, paramsJSON string) ToolResult {
	var p struct {
		Amount float64 `json:"amount"`
		Reason string  `json:"reason"`
	}
	json.Unmarshal([]byte(paramsJSON), &p)

	return ToolResult{
		Success: true,
		Message: fmt.Sprintf("Refund of $%.2f processed successfully. Transaction ID: RF-%d", p.Amount, time.Now().UnixMilli()%100000),
		Data:    map[string]any{"amount": p.Amount, "transaction_id": fmt.Sprintf("RF-%d", time.Now().UnixMilli()%100000)},
	}
}

func executeChangePlan(ctx context.Context, ticketID, paramsJSON string) ToolResult {
	var p struct {
		NewPlan string `json:"new_plan"`
	}
	json.Unmarshal([]byte(paramsJSON), &p)

	ticket, err := postgre.GetTicket(ctx, ticketID)
	if err != nil {
		return ToolResult{Success: false, Message: "ticket not found"}
	}

	customer, err := postgre.GetCustomer(ctx, ticket.CustomerID)
	if err != nil {
		return ToolResult{Success: false, Message: "customer not found"}
	}

	return ToolResult{
		Success: true,
		Message: fmt.Sprintf("Plan changed from '%s' to '%s' for customer %s", customer.Plan, p.NewPlan, customer.Name),
	}
}

func executeResetPassword(ctx context.Context, ticketID, paramsJSON string) ToolResult {
	return ToolResult{
		Success: true,
		Message: "Password reset link sent to customer's email",
	}
}

func executeEscalate(ctx context.Context, ticketID, paramsJSON string) ToolResult {
	var p struct {
		Reason   string `json:"reason"`
		Priority string `json:"priority"`
	}
	json.Unmarshal([]byte(paramsJSON), &p)

	priority := p.Priority
	if priority == "" {
		priority = constants.PriorityHigh
	}

	postgre.UpdateTicketStatus(ctx, ticketID, constants.TicketStatusWaiting)

	return ToolResult{
		Success: true,
		Message: fmt.Sprintf("Ticket escalated to senior support with priority '%s'. Reason: %s", priority, p.Reason),
	}
}

func executeSendEmail(ctx context.Context, ticketID, paramsJSON string) ToolResult {
	var p struct {
		Subject string `json:"subject"`
		Body    string `json:"body"`
	}
	json.Unmarshal([]byte(paramsJSON), &p)

	return ToolResult{
		Success: true,
		Message: fmt.Sprintf("Email sent with subject: '%s'", p.Subject),
	}
}

func executeCancelSub(ctx context.Context, ticketID, paramsJSON string) ToolResult {
	return ToolResult{
		Success: true,
		Message: "Subscription cancelled. Access remains until end of billing period",
	}
}

func executeLookupBilling(ctx context.Context, ticketID, paramsJSON string) ToolResult {
	ticket, err := postgre.GetTicket(ctx, ticketID)
	if err != nil {
		return ToolResult{Success: false, Message: "ticket not found"}
	}
	customer, _ := postgre.GetCustomer(ctx, ticket.CustomerID)

	billing := map[string]any{
		"customer":       customer.Name,
		"plan":           customer.Plan,
		"last_payment":   "$9.99 on 2026-03-01",
		"next_billing":   "2026-04-01",
		"payment_method": "Visa **** 4242",
		"payments": []map[string]any{
			{"date": "2026-03-01", "amount": 9.99, "status": "completed"},
			{"date": "2026-02-01", "amount": 9.99, "status": "completed"},
			{"date": "2026-01-15", "amount": 9.99, "status": "double_charge"},
			{"date": "2026-01-15", "amount": 9.99, "status": "double_charge"},
		},
	}

	return ToolResult{
		Success: true,
		Message: fmt.Sprintf("Billing info for %s retrieved", customer.Name),
		Data:    billing,
	}
}

func executeLookupCustomer(ctx context.Context, ticketID, paramsJSON string) ToolResult {
	ticket, err := postgre.GetTicket(ctx, ticketID)
	if err != nil {
		return ToolResult{Success: false, Message: "ticket not found"}
	}
	customer, err := postgre.GetCustomer(ctx, ticket.CustomerID)
	if err != nil {
		return ToolResult{Success: false, Message: "customer not found"}
	}

	return ToolResult{
		Success: true,
		Message: "Customer info retrieved",
		Data: map[string]any{
			"id":         customer.ID,
			"name":       customer.Name,
			"email":      customer.Email,
			"plan":       customer.Plan,
			"created_at": customer.CreatedAt,
			"tickets_count": 3,
			"lifetime_value": "$149.85",
		},
	}
}
