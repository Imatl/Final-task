package chat

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"supportflow/core/constants"
	"supportflow/core/structs"
	"supportflow/db/postgre"
	"supportflow/services/ai"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		var req structs.ChatRequest
		if err := json.Unmarshal(msg, &req); err != nil {
			sendWSError(conn, "invalid request format")
			continue
		}

		ctx := r.Context()

		if req.TicketID == "" {
			ticket := &structs.Ticket{
				CustomerID: req.CustomerID,
				Subject:    truncate(req.Message, 100),
				Status:     constants.TicketStatusOpen,
				Priority:   constants.PriorityMedium,
				Category:   "general",
			}
			if err := postgre.CreateTicket(ctx, ticket); err != nil {
				sendWSError(conn, "failed to create ticket")
				continue
			}
			req.TicketID = ticket.ID
		}

		custMsg := &structs.Message{
			TicketID: req.TicketID,
			Role:     constants.RoleCustomer,
			Content:  req.Message,
		}
		postgre.CreateMessage(ctx, custMsg)

		analysis, _ := ai.AnalyzeTicket(ctx, req.TicketID, req.Message)

		if analysis != nil {
			ticket, _ := postgre.GetTicket(ctx, req.TicketID)
			if ticket != nil {
				ticket.Priority = mapUrgencyToPriority(analysis.Urgency)
				ticket.Category = analysis.Intent
			}
		}

		sendWSJSON(conn, map[string]any{
			"type":     "analysis",
			"ticket_id": req.TicketID,
			"analysis": analysis,
		})

		resp, err := ai.ProcessMessage(ctx, req.TicketID, req.Message)
		if err != nil {
			sendWSError(conn, "AI processing failed: "+err.Error())
			continue
		}

		aiMsg := &structs.Message{
			TicketID: req.TicketID,
			Role:     constants.RoleAI,
			Content:  resp.Message,
		}
		postgre.CreateMessage(ctx, aiMsg)

		sendWSJSON(conn, map[string]any{
			"type":     "response",
			"ticket_id": req.TicketID,
			"message":  resp.Message,
			"actions":  resp.Actions,
			"auto_fixed": resp.AutoFixed,
		})
	}
}

func HandleChatHTTP(w http.ResponseWriter, r *http.Request) {
	var req structs.ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	if req.TicketID == "" {
		ticket := &structs.Ticket{
			CustomerID: req.CustomerID,
			Subject:    truncate(req.Message, 100),
			Channel:    req.Channel,
			Status:     constants.TicketStatusOpen,
			Priority:   constants.PriorityMedium,
			Category:   "general",
		}
		if err := postgre.CreateTicket(ctx, ticket); err != nil {
			http.Error(w, `{"error":"failed to create ticket"}`, http.StatusInternalServerError)
			return
		}
		req.TicketID = ticket.ID
	}

	custMsg := &structs.Message{
		TicketID: req.TicketID,
		Role:     constants.RoleCustomer,
		Content:  req.Message,
	}
	postgre.CreateMessage(ctx, custMsg)

	resp, err := ai.ProcessMessage(ctx, req.TicketID, req.Message)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	aiMsg := &structs.Message{
		TicketID: req.TicketID,
		Role:     constants.RoleAI,
		Content:  resp.Message,
	}
	postgre.CreateMessage(ctx, aiMsg)

	analysis, _ := ai.AnalyzeTicket(ctx, req.TicketID, req.Message)
	resp.Analysis = analysis

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func sendWSJSON(conn *websocket.Conn, data any) {
	msg, _ := json.Marshal(data)
	conn.WriteMessage(websocket.TextMessage, msg)
}

func sendWSError(conn *websocket.Conn, errMsg string) {
	sendWSJSON(conn, map[string]string{"type": "error", "message": errMsg})
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

func mapUrgencyToPriority(urgency string) string {
	switch urgency {
	case constants.UrgencyHigh:
		return constants.PriorityHigh
	case constants.UrgencyLow:
		return constants.PriorityLow
	default:
		return constants.PriorityMedium
	}
}
