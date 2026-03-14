package analytics

import (
	"encoding/json"
	"log"
	"net/http"

	"supportflow/db/postgre"
)

func HandleOverview(w http.ResponseWriter, r *http.Request) {
	company := r.URL.Query().Get("company")

	overview, err := postgre.GetAnalyticsOverview(r.Context(), company)
	if err != nil {
		log.Printf("[analytics] overview error: %v", err)
		http.Error(w, `{"error":"failed to fetch analytics"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(overview)
}

func HandleAgentPerformance(w http.ResponseWriter, r *http.Request) {
	company := r.URL.Query().Get("company")

	perfs, err := postgre.GetAgentPerformance(r.Context(), company)
	if err != nil {
		log.Printf("[analytics] agent performance error: %v", err)
		http.Error(w, `{"error":"failed to fetch performance"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(perfs)
}
