package routes

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"

	"supportflow/api/analytics"
	"supportflow/api/auth"
	"supportflow/api/chat"
	"supportflow/api/companies"
	apiIntegrations "supportflow/api/integrations"
	"supportflow/api/knowledge"
	"supportflow/api/settings"
	"supportflow/api/tickets"
)

func Register(r *mux.Router, ctx context.Context) {
	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	}).Methods("GET")

	api.HandleFunc("/auth/google", auth.HandleGoogleAuth).Methods("GET")
	api.HandleFunc("/auth/google/callback", auth.HandleGoogleCallback).Methods("GET")
	api.HandleFunc("/auth/login", auth.HandleLogin).Methods("POST")
	api.HandleFunc("/auth/register", auth.HandleRegister).Methods("POST")
	api.HandleFunc("/auth/invite", auth.HandleGenerateInvite).Methods("POST")
	api.HandleFunc("/invite/{token}", auth.HandleValidateInvite).Methods("GET")

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

	api.HandleFunc("/companies", companies.HandleList).Methods("GET")

	api.HandleFunc("/integrations", apiIntegrations.HandleList).Methods("GET")
	api.HandleFunc("/integrations/connect", apiIntegrations.HandleConnect).Methods("POST")
	api.HandleFunc("/integrations/disconnect", apiIntegrations.HandleDisconnect).Methods("POST")

	api.HandleFunc("/settings/providers", settings.HandleGetProviders).Methods("GET")
	api.HandleFunc("/settings/providers", settings.HandleSetProvider).Methods("PUT")
	api.HandleFunc("/settings/metrics", settings.HandleGetMetrics).Methods("GET")

	api.HandleFunc("/knowledge", knowledge.HandleList).Methods("GET")
	api.HandleFunc("/knowledge", knowledge.HandleCreate).Methods("POST")
	api.HandleFunc("/knowledge/{id}", knowledge.HandleDelete).Methods("DELETE")
}
