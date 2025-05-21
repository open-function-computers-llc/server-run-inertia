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
			if acc.Name == r.PathValue("accountName") {
				a = acc
				foundAccount = true
				break
			}
		}

		if !foundAccount {
			s.inertiaManager.Render(w, r, "Error", map[string]any{
				"account": a, // TODO: fix this prob bundle
			})
			return
		}

		tabMapper := map[string]string{
			"":          "Account/Show",
			"domains":   "Account/Domains",
			"analytics": "Account/Analytics",
		}
		s.inertiaManager.Render(w, r, tabMapper[tab], map[string]any{
			"account": a,
		})
	}
}
