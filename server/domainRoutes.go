package server

import (
	"encoding/json"
	"net/http"
)

func (s *server) domainRoutes() map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{
		"POST /accounts/{account}/domains":            s.handleAddDomain(),
		"PUT /accounts/{account}/domains/primary":     s.handleSetPrimaryDomain(),
		"DELETE /accounts/{account}/domains/{domain}": s.handleDeleteDomain(),
	}
}

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

func (s *server) handleSetPrimaryDomain() http.HandlerFunc {
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

func (s *server) handleDeleteDomain() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// accountName := r.PathValue("account")
		domain := r.PathValue("domain")

		if domain == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Domain is required"})
			return
		}

		// TODO: Add your database logic here
		
		w.WriteHeader(http.StatusOK)
	}
}
