package server

import (
	"encoding/json"
	"net/http"
)

func (s *server) handleFormSetUptimeURL() http.HandlerFunc {
	type incomingPayload struct {
		UptimeURL string `json:"uptimeURL"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		a, err := s.findAccount(r.PathValue("account"))
		if err != nil {
			s.flashError(w, err.Error())
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			return
		}

		var payload incomingPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			s.flashError(w, "Problem reading payload")
			http.Redirect(w, r, "/account/"+a.Name, http.StatusSeeOther)
			return
		}

		if payload.UptimeURL == "" {
			s.flashError(w, "Uptime URL is required")
			http.Redirect(w, r, "/account/"+a.Name, http.StatusSeeOther)
			return
		}

		if err := a.SetUptimeURL(payload.UptimeURL); err != nil {
			s.flashError(w, "Couldn't save the uptime URL: "+err.Error())
		}

		http.Redirect(w, r, "/account/"+a.Name, http.StatusSeeOther)
	}
}
