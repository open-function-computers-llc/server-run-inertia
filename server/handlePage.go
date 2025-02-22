package server

import (
	"net/http"
)

func (s *server) handlePage(pageName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if pageName == "Index" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		s.inertiaManager.Render(w, r, pageName, nil)
	}
}
