package server

import (
	"net/http"
)

func (s *server) handleSettings() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.inertiaManager.Render(w, r, "Settings", map[string]any{})
	}
}
