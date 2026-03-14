package tickets

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"supportflow/core/structs"
	"supportflow/db/postgre"
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
		http.Error(w, `{"error":"ticket not found"}`, http.StatusNotFound)
		return
	}

	customer, _ := postgre.GetCustomer(r.Context(), ticket.CustomerID)
	messages, _ := postgre.GetMessagesByTicket(r.Context(), id)
	actions, _ := postgre.GetActionsByTicket(r.Context(), id)

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
		http.Error(w, `{"error":"update failed"}`, http.StatusInternalServerError)
		return
	}

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
		http.Error(w, `{"error":"assign failed"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"ok":true}`))
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
		http.Error(w, `{"error":"update failed"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"ok":true}`))
}
