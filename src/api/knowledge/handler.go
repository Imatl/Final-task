package knowledge

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"supportflow/db/postgre"
)

func HandleList(w http.ResponseWriter, r *http.Request) {
	company := r.URL.Query().Get("company")
	if company == "" {
		http.Error(w, `{"error":"company is required"}`, http.StatusBadRequest)
		return
	}

	entries, err := postgre.ListKBEntries(r.Context(), company)
	if err != nil {
		log.Printf("[knowledge] list error: %v", err)
		http.Error(w, `{"error":"failed to list entries"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entries)
}

func HandleCreate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Company  string `json:"company"`
		Question string `json:"question"`
		Answer   string `json:"answer"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}

	if req.Company == "" || req.Question == "" || req.Answer == "" {
		http.Error(w, `{"error":"company, question and answer are required"}`, http.StatusBadRequest)
		return
	}

	entry, err := postgre.CreateKBEntry(r.Context(), req.Company, req.Question, req.Answer)
	if err != nil {
		log.Printf("[knowledge] create error: %v", err)
		http.Error(w, `{"error":"failed to create entry"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entry)
}

func HandleDelete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		http.Error(w, `{"error":"id is required"}`, http.StatusBadRequest)
		return
	}

	if err := postgre.DeleteKBEntry(r.Context(), id); err != nil {
		log.Printf("[knowledge] delete error: %v", err)
		http.Error(w, `{"error":"failed to delete entry"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"ok":true}`))
}
