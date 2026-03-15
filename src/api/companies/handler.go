package companies

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"supportflow/db/postgre"
)

type CompanyInfo struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	StaffCount  int       `json:"staff_count"`
	AISpendUSD  float64   `json:"ai_spend_usd"`
	CreatedAt   time.Time `json:"created_at"`
}

func HandleList(w http.ResponseWriter, r *http.Request) {
	companies, err := postgre.ListCompanies(r.Context())
	if err != nil {
		log.Printf("[companies] list error: %v", err)
		http.Error(w, `{"error":"failed to list companies"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"companies": companies,
		"total":     len(companies),
	})
}
