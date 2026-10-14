package server

import (
	"encoding/json"
	"net/http"
)

func (s *server) handleVcsAddDeployKey() http.HandlerFunc {
	type incomingPayload struct {
		Repo     string `json:"repo"`
		Name     string `json:"name"`
		Key      string `json:"key"`
		ReadOnly bool   `json:"readOnly"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var payload incomingPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			s.flashError(w, "Problem reading payload")
			http.Redirect(w, r, "/vcs", http.StatusSeeOther)
			return
		}

		if payload.Repo == "" || payload.Name == "" || payload.Key == "" {
			s.flashError(w, "Repo, name, and key are all required")
			http.Redirect(w, r, "/vcs", http.StatusSeeOther)
			return
		}

		if err := s.vcsProvider.AddDeployKey(payload.Repo, payload.Name, payload.Key, payload.ReadOnly); err != nil {
			s.flashError(w, err.Error())
			http.Redirect(w, r, "/vcs", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/vcs", http.StatusSeeOther)
	}
}
