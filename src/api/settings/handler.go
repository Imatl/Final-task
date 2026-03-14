package settings

import (
	"encoding/json"
	"net/http"
	"strconv"

	"supportflow/services/ai"
)

func HandleGetProviders(w http.ResponseWriter, r *http.Request) {
	_, active := ai.GetActiveProvider()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"providers": ai.GetProviderNames(),
		"active":    active,
	})
}

func HandleSetProvider(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Provider string `json:"provider"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Provider == "" {
		http.Error(w, `{"error":"provider required"}`, http.StatusBadRequest)
		return
	}

	if !ai.SetActiveProvider(body.Provider) {
		http.Error(w, `{"error":"unknown provider"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"ok":true}`))
}

func HandleGetMetrics(w http.ResponseWriter, r *http.Request) {
	limit := 100
	if l, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil {
		limit = l
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"metrics": ai.GetMetrics(limit),
		"stats":   ai.GetMetricsStats(),
	})
}
