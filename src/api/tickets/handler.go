package tickets

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"supportflow/core/structs"
	"supportflow/db/postgre"
	"supportflow/services/ai"
)

func HandleList(w http.ResponseWriter, r *http.Request) {
	filter := structs.TicketFilter{
		Status:   r.URL.Query().Get("status"),
		Priority: r.URL.Query().Get("priority"),
		AgentID:  r.URL.Query().Get("agent_id"),
		Category: r.URL.Query().Get("category"),
	}
	if l, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil {
		filter.Limit = l
	}
	if o, err := strconv.Atoi(r.URL.Query().Get("offset")); err == nil {
		filter.Offset = o
	}

	tickets, total, err := postgre.ListTickets(r.Context(), filter)
	if err != nil {
		log.Printf("[tickets] list error: %v", err)
		http.Error(w, `{"error":"failed to fetch tickets"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(structs.TicketListResponse{Tickets: tickets, Total: total})
}

func HandleGet(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	ticket, err := postgre.GetTicket(r.Context(), id)
	if err != nil {
		log.Printf("[tickets] get %s error: %v", id, err)
		http.Error(w, `{"error":"ticket not found"}`, http.StatusNotFound)
		return
	}

	customer, err := postgre.GetCustomer(r.Context(), ticket.CustomerID)
	if err != nil {
		log.Printf("[tickets] get customer %s error: %v", ticket.CustomerID, err)
	}
	messages, err := postgre.GetMessagesByTicket(r.Context(), id)
	if err != nil {
		log.Printf("[tickets] get messages for %s error: %v", id, err)
	}
	actions, err := postgre.GetActionsByTicket(r.Context(), id)
	if err != nil {
		log.Printf("[tickets] get actions for %s error: %v", id, err)
	}

	detail := structs.TicketDetail{
		Ticket:   *ticket,
		Messages: messages,
		Actions:  actions,
	}
	if customer != nil {
		detail.Customer = *customer
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(detail)
}

func HandleUpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var body struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Status == "" {
		http.Error(w, `{"error":"status required"}`, http.StatusBadRequest)
		return
	}

	if err := postgre.UpdateTicketStatus(r.Context(), id, body.Status); err != nil {
		log.Printf("[tickets] update status %s -> %s error: %v", id, body.Status, err)
		http.Error(w, `{"error":"update failed"}`, http.StatusInternalServerError)
		return
	}

	log.Printf("[tickets] status updated %s -> %s", id, body.Status)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"ok":true}`))
}

func HandleAssign(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var body struct {
		AgentID string `json:"agent_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.AgentID == "" {
		http.Error(w, `{"error":"agent_id required"}`, http.StatusBadRequest)
		return
	}

	if err := postgre.AssignTicket(r.Context(), id, body.AgentID); err != nil {
		log.Printf("[tickets] assign %s -> agent %s error: %v", id, body.AgentID, err)
		http.Error(w, `{"error":"assign failed"}`, http.StatusInternalServerError)
		return
	}

	log.Printf("[tickets] assigned %s -> agent %s", id, body.AgentID)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"ok":true}`))
}

func HandleAgentReply(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var body struct {
		AgentID string `json:"agent_id"`
		Message string `json:"message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Message == "" {
		http.Error(w, `{"error":"message required"}`, http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	if _, err := postgre.GetTicket(ctx, id); err != nil {
		log.Printf("[tickets] agent reply - ticket %s not found: %v", id, err)
		http.Error(w, `{"error":"ticket not found"}`, http.StatusNotFound)
		return
	}

	msg := &structs.Message{
		TicketID: id,
		Role:     "agent",
		Content:  body.Message,
	}
	if err := postgre.CreateMessage(ctx, msg); err != nil {
		log.Printf("[tickets] save agent message error: %v", err)
		http.Error(w, `{"error":"failed to save message"}`, http.StatusInternalServerError)
		return
	}

	if body.AgentID != "" {
		_ = postgre.AssignTicket(ctx, id, body.AgentID)
	}
	_ = postgre.UpdateTicketStatus(ctx, id, "in_progress")

	log.Printf("[tickets] agent reply on %s by %s", id, body.AgentID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"ok": true, "message": msg})
}

func HandleSuggest(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	suggestion, err := ai.GenerateSuggestion(r.Context(), id)
	if err != nil {
		log.Printf("[tickets] suggest for %s error: %v", id, err)
		http.Error(w, `{"error":"failed to generate suggestion"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"suggestion": suggestion})
}

func HandleApproveAction(w http.ResponseWriter, r *http.Request) {
	var body structs.ActionApproval
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}

	status := "approved"
	if !body.Approved {
		status = "rejected"
	}

	if err := postgre.UpdateActionStatus(r.Context(), body.ActionID, status, ""); err != nil {
		log.Printf("[tickets] approve action %s error: %v", body.ActionID, err)
		http.Error(w, `{"error":"update failed"}`, http.StatusInternalServerError)
		return
	}

	log.Printf("[tickets] action %s -> %s", body.ActionID, status)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"ok":true}`))
}
