package routes

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"

	"supportflow/api/analytics"
	"supportflow/api/chat"
	"supportflow/api/settings"
	"supportflow/api/tickets"
)

func Register(r *mux.Router, ctx context.Context) {
	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	}).Methods("GET")

	api.HandleFunc("/chat", chat.HandleChatHTTP).Methods("POST")
	api.HandleFunc("/chat/ws", chat.HandleWebSocket)

	api.HandleFunc("/tickets", tickets.HandleList).Methods("GET")
	api.HandleFunc("/tickets/{id}", tickets.HandleGet).Methods("GET")
	api.HandleFunc("/tickets/{id}/status", tickets.HandleUpdateStatus).Methods("PUT")
	api.HandleFunc("/tickets/{id}/assign", tickets.HandleAssign).Methods("PUT")
	api.HandleFunc("/tickets/{id}/reply", tickets.HandleAgentReply).Methods("POST")
	api.HandleFunc("/tickets/{id}/suggest", tickets.HandleSuggest).Methods("POST")
	api.HandleFunc("/tickets/actions/approve", tickets.HandleApproveAction).Methods("POST")

	api.HandleFunc("/analytics/overview", analytics.HandleOverview).Methods("GET")
	api.HandleFunc("/analytics/agents", analytics.HandleAgentPerformance).Methods("GET")

	api.HandleFunc("/settings/providers", settings.HandleGetProviders).Methods("GET")
	api.HandleFunc("/settings/providers", settings.HandleSetProvider).Methods("PUT")
	api.HandleFunc("/settings/metrics", settings.HandleGetMetrics).Methods("GET")
}
