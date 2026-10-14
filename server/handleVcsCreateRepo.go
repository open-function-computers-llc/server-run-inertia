package server

import (
	"encoding/json"
	"net/http"
)

func (s *server) handleVcsCreateRepo() http.HandlerFunc {
	type incomingPayload struct {
		Name string `json:"name"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var payload incomingPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			s.flashError(w, "Problem reading payload")
			http.Redirect(w, r, "/vcs", http.StatusSeeOther)
			return
		}

		if payload.Name == "" {
			s.flashError(w, "Repository name is required")
			http.Redirect(w, r, "/vcs", http.StatusSeeOther)
			return
		}

		if err := s.vcsProvider.CreateRepository(payload.Name); err != nil {
			s.flashError(w, err.Error())
			http.Redirect(w, r, "/vcs", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/vcs", http.StatusSeeOther)
	}
}
