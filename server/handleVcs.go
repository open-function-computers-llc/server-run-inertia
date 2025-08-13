package server

import (
	"net/http"
)

func (s *server) handleVcs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		repos, _ := s.vcsProvider.ListRepositories()

		s.inertiaManager.Render(w, r, "VCS", map[string]any{
			"error": s.vcsErrorMessage,
			"repos": repos,
		})
	}
}
