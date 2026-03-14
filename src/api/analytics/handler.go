package analytics

import (
	"encoding/json"
	"net/http"

	"supportflow/db/postgre"
)

func HandleOverview(w http.ResponseWriter, r *http.Request) {
	overview, err := postgre.GetAnalyticsOverview(r.Context())
	if err != nil {
		http.Error(w, `{"error":"failed to fetch analytics"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(overview)
}

func HandleAgentPerformance(w http.ResponseWriter, r *http.Request) {
	perfs, err := postgre.GetAgentPerformance(r.Context())
	if err != nil {
		http.Error(w, `{"error":"failed to fetch performance"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(perfs)
}
