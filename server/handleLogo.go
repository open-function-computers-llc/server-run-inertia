package server

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"
)

func (s *server) handleLogo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// default to the baked in OFC logo
		logoBytes := s.defaultLogo

		// check to see if we should load up an alternate logo
		logoFile := os.Getenv("LOGO_PNG")
		if logoFile != "" {
			logoBytes, _ = os.ReadFile(logoFile)
		}

		output := map[string]string{
			"value": base64.StdEncoding.EncodeToString(logoBytes),
		}
		bytes, _ := json.Marshal(output)
		w.Write(bytes)
	}
}
