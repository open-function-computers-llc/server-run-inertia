package server

import "net/http"

func (s *server) handleHome() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.inertiaManager.Render(w, r, "Home", nil)
	}
}