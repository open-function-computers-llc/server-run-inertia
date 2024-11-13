package server

import (
	"encoding/json"
	"io"
	"net/http"
)

func (s *server) handleFormAccountCreate() http.HandlerFunc {
	type incomingPayload struct {
		NewAccountName string `json:"newAccountName"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var payload incomingPayload
		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			sendJSON(w, "Error reading payload body")
			return
		}
		err = json.Unmarshal(bytes, &payload)
		if err != nil {
			sendJSON(w, "Invalid JSON Body")
			return
		}

		if payload.NewAccountName == "" {
			s.inertiaManager.Share("error", "newAccountName is required")
			http.Redirect(w, r, "/accounts", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/account/"+payload.NewAccountName, http.StatusSeeOther)
	}
}
