package server

import (
	"net/http"
)

func (s *server) handleVcs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vcsError := s.vcsErrorMessage
		repos := []string{}

		if s.vcsAvailable {
			var err error
			repos, err = s.vcsProvider.ListRepositories()
			if err != nil {
				vcsError = "Couldn't list repositories: " + err.Error()
			}
		} else if vcsError == "" {
			vcsError = "Version control isn't configured on this server. Set VCS_PROVIDER to GITEA or GITHUB."
		}

		s.render(w, r, "VCS", map[string]any{
			"vcsError":     vcsError,
			"repos":        repos,
			"vcsAvailable": s.vcsAvailable,
		})
	}
}
