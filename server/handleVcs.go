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
			vcsError = "Version control isn't configured on this server. Set GITEA_URL and GITEA_TOKEN (default provider) or GITHUB_TOKEN, and optionally VCS_PROVIDER to change the default."
		}

		s.render(w, r, "VCS", map[string]any{
			"vcsError":     vcsError,
			"repos":        repos,
			"vcsAvailable": s.vcsAvailable,
		})
	}
}
