package server

import (
	"net/http"

	"github.com/open-function-computers-llc/server-run-inertia/session"
)

func (s *server) handleLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie("server-run-auth")
		if err != nil {
			ctx := s.inertiaManager.WithProp(r.Context(), "errMessage", "no cookie")
			http.Redirect(w, r.WithContext(ctx), "/login", http.StatusSeeOther)
			return
		}

		session.Invalidate(sessionCookie.Value)

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
