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
		log.Printf("[ws] upgrade error: %v", err)
		return
	}
	defer conn.Close()
	log.Println("[ws] client connected")

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("[ws] read error (client disconnected): %v", err)
			break
		}

		var req structs.ChatRequest
		if err := json.Unmarshal(msg, &req); err != nil {
			log.Printf("[ws] invalid request: %v", err)
			sendWSError(conn, "invalid request format")
			continue
		}

		log.Printf("[ws] message from customer=%s: %s", req.CustomerID, truncate(req.Message, 50))
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
				log.Printf("[ws] create ticket error: %v", err)
				sendWSError(conn, "failed to create ticket")
				continue
			}
			req.TicketID = ticket.ID
			log.Printf("[ws] created ticket %s", ticket.ID)
		}

		custMsg := &structs.Message{
			TicketID: req.TicketID,
			Role:     constants.RoleCustomer,
			Content:  req.Message,
		}
		if err := postgre.CreateMessage(ctx, custMsg); err != nil {
			log.Printf("[ws] save customer message error: %v", err)
		}

		analysis, err := ai.AnalyzeTicket(ctx, req.TicketID, req.Message)
		if err != nil {
			log.Printf("[ws] analyze ticket error: %v", err)
		}

		if analysis != nil {
			ticket, err := postgre.GetTicket(ctx, req.TicketID)
			if err != nil {
				log.Printf("[ws] get ticket error: %v", err)
			}
			if ticket != nil {
				ticket.Priority = mapUrgencyToPriority(analysis.Urgency)
				ticket.Category = analysis.Intent
			}
		}

		sendWSJSON(conn, map[string]any{
			"type":      "analysis",
			"ticket_id": req.TicketID,
			"analysis":  analysis,
		})

		resp, err := ai.ProcessMessage(ctx, req.TicketID, req.Message, "")
		if err != nil {
			log.Printf("[ws] AI process error: %v", err)
			sendWSError(conn, "AI processing failed: "+err.Error())
			continue
		}

		aiMsg := &structs.Message{
			TicketID: req.TicketID,
			Role:     constants.RoleAI,
			Content:  resp.Message,
		}
		if err := postgre.CreateMessage(ctx, aiMsg); err != nil {
			log.Printf("[ws] save AI message error: %v", err)
		}

		sendWSJSON(conn, map[string]any{
			"type":       "response",
			"ticket_id":  req.TicketID,
			"message":    resp.Message,
			"actions":    resp.Actions,
			"auto_fixed": resp.AutoFixed,
		})
	}
}

func HandleChatHTTP(w http.ResponseWriter, r *http.Request) {
	var req structs.ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[chat] decode request error: %v", err)
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}

	log.Printf("[chat] message from customer=%s channel=%s: %s", req.CustomerID, req.Channel, truncate(req.Message, 50))
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
			log.Printf("[chat] create ticket error: %v", err)
			http.Error(w, `{"error":"failed to create ticket"}`, http.StatusInternalServerError)
			return
		}
		req.TicketID = ticket.ID
		log.Printf("[chat] created ticket %s", ticket.ID)
	}

	custMsg := &structs.Message{
		TicketID: req.TicketID,
		Role:     constants.RoleCustomer,
		Content:  req.Message,
	}
	if err := postgre.CreateMessage(ctx, custMsg); err != nil {
		log.Printf("[chat] save customer message error: %v", err)
	}

	resp, err := ai.ProcessMessage(ctx, req.TicketID, req.Message, "")
	if err != nil {
		log.Printf("[chat] AI error: %v", err)
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	aiMsg := &structs.Message{
		TicketID: req.TicketID,
		Role:     constants.RoleAI,
		Content:  resp.Message,
	}
	if err := postgre.CreateMessage(ctx, aiMsg); err != nil {
		log.Printf("[chat] save AI message error: %v", err)
	}

	analysis, err := ai.AnalyzeTicket(ctx, req.TicketID, req.Message)
	if err != nil {
		log.Printf("[chat] analyze error: %v", err)
	}
	resp.Analysis = analysis

	if analysis != nil && analysis.Sentiment == constants.SentimentPositive && !resp.AutoFixed {
		actions, _ := postgre.GetActionsByTicket(ctx, req.TicketID)
		if len(actions) > 0 {
			if err := postgre.UpdateTicketStatus(ctx, req.TicketID, constants.TicketStatusResolved); err != nil {
				log.Printf("[chat] auto-resolve on positive sentiment error: %v", err)
			} else {
				log.Printf("[chat] ticket %s auto-resolved (positive sentiment after actions)", req.TicketID)
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func sendWSJSON(conn *websocket.Conn, data any) {
	msg, err := json.Marshal(data)
	if err != nil {
		log.Printf("[ws] marshal error: %v", err)
		return
	}
	if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
		log.Printf("[ws] write error: %v", err)
	}
}

func sendWSError(conn *websocket.Conn, errMsg string) {
	log.Printf("[ws] sending error to client: %s", errMsg)
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
