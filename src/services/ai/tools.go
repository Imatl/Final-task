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
	log.Printf("[tool] executing %s for ticket %s params=%s", toolName, ticketID, paramsJSON)

	var result ToolResult
	switch toolName {
	case constants.ActionRefund:
		result = executeRefund(ctx, ticketID, paramsJSON)
	case constants.ActionChangePlan:
		result = executeChangePlan(ctx, ticketID, paramsJSON)
	case constants.ActionResetPassword:
		result = executeResetPassword(ctx, ticketID, paramsJSON)
	case constants.ActionEscalate:
		result = executeEscalate(ctx, ticketID, paramsJSON)
	case constants.ActionSendEmail:
		result = executeSendEmail(ctx, ticketID, paramsJSON)
	case constants.ActionCancelSub:
		result = executeCancelSub(ctx, ticketID, paramsJSON)
	case "lookup_billing":
		result = executeLookupBilling(ctx, ticketID, paramsJSON)
	case "lookup_customer":
		result = executeLookupCustomer(ctx, ticketID, paramsJSON)
	default:
		result = ToolResult{Success: false, Message: fmt.Sprintf("unknown tool: %s", toolName)}
	}

	log.Printf("[tool] %s result: success=%v message=%s", toolName, result.Success, result.Message)
	return result
}

func executeRefund(ctx context.Context, ticketID, paramsJSON string) ToolResult {
	var p struct {
		Amount float64 `json:"amount"`
		Reason string  `json:"reason"`
	}
	if err := json.Unmarshal([]byte(paramsJSON), &p); err != nil {
		log.Printf("[tool] refund unmarshal error: %v", err)
	}

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
	if err := json.Unmarshal([]byte(paramsJSON), &p); err != nil {
		log.Printf("[tool] change_plan unmarshal error: %v", err)
	}

	ticket, err := postgre.GetTicket(ctx, ticketID)
	if err != nil {
		log.Printf("[tool] change_plan get ticket error: %v", err)
		return ToolResult{Success: false, Message: "ticket not found"}
	}

	customer, err := postgre.GetCustomer(ctx, ticket.CustomerID)
	if err != nil {
		log.Printf("[tool] change_plan get customer error: %v", err)
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
	if err := json.Unmarshal([]byte(paramsJSON), &p); err != nil {
		log.Printf("[tool] escalate unmarshal error: %v", err)
	}

	priority := p.Priority
	if priority == "" {
		priority = constants.PriorityHigh
	}

	if err := postgre.UpdateTicketStatus(ctx, ticketID, constants.TicketStatusWaiting); err != nil {
		log.Printf("[tool] escalate update status error: %v", err)
		return ToolResult{
			Success: false,
			Message: fmt.Sprintf("failed to escalate: %v", err),
		}
	}

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
	if err := json.Unmarshal([]byte(paramsJSON), &p); err != nil {
		log.Printf("[tool] send_email unmarshal error: %v", err)
	}

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

var planPrices = map[string]float64{
	"free":    0,
	"basic":   9.99,
	"premium": 29.99,
}

func executeLookupBilling(ctx context.Context, ticketID, paramsJSON string) ToolResult {
	ticket, err := postgre.GetTicket(ctx, ticketID)
	if err != nil {
		log.Printf("[tool] lookup_billing get ticket error: %v", err)
		return ToolResult{Success: false, Message: "ticket not found"}
	}
	customer, err := postgre.GetCustomer(ctx, ticket.CustomerID)
	if err != nil {
		log.Printf("[tool] lookup_billing get customer error: %v", err)
		return ToolResult{Success: false, Message: "customer not found"}
	}

	price := planPrices[customer.Plan]
	now := time.Now()

	var payments []any
	if price > 0 {
		for i := 0; i < 3; i++ {
			d := now.AddDate(0, -i, 0)
			payments = append(payments, map[string]any{
				"date":   d.Format("2006-01-02"),
				"amount": price,
				"status": "paid",
				"method": "card ending 4242",
			})
		}
	}

	data := map[string]any{
		"customer":        customer.Name,
		"email":           customer.Email,
		"plan":            customer.Plan,
		"plan_price":      price,
		"billing_status":  "active",
		"recent_payments": payments,
		"open_issues":     []any{},
	}
	if customer.Phone != nil {
		data["phone"] = *customer.Phone
	}

	return ToolResult{
		Success: true,
		Message: fmt.Sprintf("Billing info for %s retrieved", customer.Name),
		Data:    data,
	}
}

func executeLookupCustomer(ctx context.Context, ticketID, paramsJSON string) ToolResult {
	ticket, err := postgre.GetTicket(ctx, ticketID)
	if err != nil {
		log.Printf("[tool] lookup_customer get ticket error: %v", err)
		return ToolResult{Success: false, Message: "ticket not found"}
	}
	customer, err := postgre.GetCustomer(ctx, ticket.CustomerID)
	if err != nil {
		log.Printf("[tool] lookup_customer get customer error: %v", err)
		return ToolResult{Success: false, Message: "customer not found"}
	}

	data := map[string]any{
		"id":         customer.ID,
		"name":       customer.Name,
		"email":      customer.Email,
		"plan":       customer.Plan,
		"created_at": customer.CreatedAt,
	}
	if customer.Phone != nil {
		data["phone"] = *customer.Phone
	}

	return ToolResult{
		Success: true,
		Message: "Customer info retrieved",
		Data:    data,
	}
}
