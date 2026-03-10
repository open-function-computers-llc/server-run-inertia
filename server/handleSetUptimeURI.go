package server

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/open-function-computers-llc/server-run-inertia/account"
)

func (s *server) handleSetUptimeURI() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountName := r.PathValue("account")

		var req struct {
			UptimeURI string `json:"uptimeURI"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request"})
			return
		}

		if strings.TrimSpace(req.UptimeURI) == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Uptime URI is required"})
			return
		}

		a, err := account.Find(accountName)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Account not found"})
			return
		}

		if err := a.SetUptimeURI(req.UptimeURI); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to save"})
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
