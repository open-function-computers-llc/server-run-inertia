package server

import (
	"net/http"
	"os"
)

func (s *server) handleCreateAccount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.render(w, r, "Account/Create", map[string]any{
			"vcsAvailable":  s.vcsAvailable,
			"stagingDomain": stagingDomain(),
		})
	}
}

// stagingDomain is where new sites are served when no domain is given
func stagingDomain() string {
	if d := os.Getenv("STAGING_DOMAIN"); d != "" {
		return d
	}
	return "stagingsubdomain.com"
}
