package server

import (
	"encoding/json"
	"net/http"

	"github.com/open-function-computers-llc/server-run-inertia/account"
)

func (s *server) handleFormSetLockedDownDisabled() http.HandlerFunc {
	type incomingPayload struct {
		LockedDownDisabled bool `json:"lockedDownDisabled"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		accountName := r.PathValue("account")
		if accountName == "" {
			s.flashError(w, "Invalid account")
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
			s.flashError(w, "Account not found: "+r.PathValue("account"))
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			return
		}

		var payload incomingPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			s.flashError(w, "Problem reading payload")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		err := a.SetLockedDownDisabled(payload.LockedDownDisabled)
		if err != nil {
			s.flashError(w, "Failed to update lock-down setting")
			http.Redirect(w, r, "/account/"+accountName, http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/account/"+accountName, http.StatusSeeOther)
	}
}
