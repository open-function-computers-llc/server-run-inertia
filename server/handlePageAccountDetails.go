package server

import (
	"net/http"

	"github.com/open-function-computers-llc/server-run-inertia/account"
)

func (s *server) handleAccountDetails(tab string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.bootstrapAccounts()

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

		tabMapper := map[string]string{
			"":          "Account/Show",
			"domains":   "Account/Domains",
			"analytics": "Account/Analytics",
			"logs":      "Account/Logs",
			"export":    "Account/Export",
			"settings":  "Account/Settings",
		}
		s.inertiaManager.Render(w, r, tabMapper[tab], map[string]any{
			"account": a,
		})
	}
}
