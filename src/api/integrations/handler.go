package integrations

import (
	"encoding/json"
	"log"
	"net/http"

	svcIntegrations "supportflow/services/integrations"
)

func HandleConnect(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ID     string            `json:"id"`
		Type   string            `json:"type"`
		Config map[string]string `json:"config"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Type == "" {
		http.Error(w, `{"error":"type and config required"}`, http.StatusBadRequest)
		return
	}

	if body.ID == "" {
		body.ID = body.Type
	}

	if err := svcIntegrations.Connect(r.Context(), body.ID, body.Type, body.Config); err != nil {
		log.Printf("[integrations] connect %s error: %v", body.Type, err)
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "connected", "id": body.ID})
}

func HandleDisconnect(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.ID == "" {
		http.Error(w, `{"error":"id required"}`, http.StatusBadRequest)
		return
	}

	if err := svcIntegrations.Disconnect(r.Context(), body.ID); err != nil {
		log.Printf("[integrations] disconnect %s error: %v", body.ID, err)
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"disconnected"}`))
}

func HandleList(w http.ResponseWriter, r *http.Request) {
	list := svcIntegrations.List()
	if list == nil {
		list = []svcIntegrations.Integration{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}
