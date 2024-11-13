package server

import "net/http"

func (s *server) handleAccounts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.bootstrapAccounts()
		s.inertiaManager.Render(w, r, "Account/Index", map[string]any{
			"accounts": s.accounts,
		})
	}
}
