package server

import (
	"net/http"
	"os"
)

func (s *server) thirdPartyProtection(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accessToken := r.Header.Get("X-Server-Run-Access")

		correctToken := os.Getenv("THIRD_PARTY_ACCESS_TOKEN")
		if accessToken != correctToken {
			sendJSONError(w, http.StatusUnauthorized, map[string]string{
				"error": "Invalid X-Server-Run-Access header provided",
			})
			return
		}

		next(w, r)
	}
}
