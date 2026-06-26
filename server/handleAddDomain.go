package server

import (
	"encoding/json"
	"net/http"
)

func (s *server) handleAddDomain() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// accountName := r.PathValue("account")

		var req struct {
			Domain string `json:"domain"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request"})
			return
		}

		if req.Domain == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"domain": "Domain is required"})
			return
		}

		// TODO: Add your database logic here

		w.WriteHeader(http.StatusOK)
	}
}
