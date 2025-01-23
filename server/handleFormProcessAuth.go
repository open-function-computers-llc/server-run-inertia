package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/open-function-computers-llc/server-run-inertia/session"
)

func (s *server) handleFormProcessAuth() http.HandlerFunc {
	type incomingPayload struct {
		Username string `json:"userName"`
		Password string `json:"password"`
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

		if payload.Password == "" || payload.Username == "" {
			s.inertiaManager.Share("error", "Both username and password are required")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		if payload.Username != s.authUser || payload.Password != s.authPassword {
			fmt.Println(payload.Username, s.authUser, payload.Password, s.authPassword)
			s.inertiaManager.Share("error", "Invalid Credentials")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		sessionKey := session.Create()

		// set a cookie on the browser!
		c := http.Cookie{
			Name:    "server-run-auth",
			Value:   sessionKey,
			Expires: time.Now().Add(time.Minute * 30),
		}
		http.SetCookie(w, &c)

		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}
}
