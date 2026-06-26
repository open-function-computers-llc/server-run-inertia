package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/open-function-computers-llc/server-run-inertia/account"
)

func (s *server) handleFormSetUptimeURL() http.HandlerFunc {
	type incomingPayload struct {
		UptimeURL string `json:"uptimeURL"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		accountName := r.PathValue("account")
		if accountName == "" {
			s.inertiaManager.Share("error", "Invalid account")
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			return
		}
		var a account.Account
		foundAccount := false
		for _, acc := range s.accounts {
			if acc.Name == r.PathValue("account") {
				a = acc
				foundAccount = true
				break
			}
		}

		if !foundAccount {
			s.inertiaManager.Share("error", "Account not found: "+r.PathValue("account"))
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			return
		}

		var payload incomingPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			s.inertiaManager.Share("error", "Problem reading payload")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		if payload.UptimeURL == "" {
			s.inertiaManager.Share("error", "Uptime URL is required")
			http.Redirect(w, r, "/account/"+a.Name, http.StatusSeeOther)
			return
		}

		fmt.Println("uptimeurl: " + payload.UptimeURL)

		err := a.SetUptimeURL(payload.UptimeURL)
		if err != nil {
			// handle couldn't update the settings json blob
		}

		http.Redirect(w, r, "/account/"+accountName, http.StatusSeeOther)
	}
}
