package server

import "net/http"

func (s *server) handleDashboard() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.inertiaManager.Render(w, r, "Dashboard", nil)
	}
}
