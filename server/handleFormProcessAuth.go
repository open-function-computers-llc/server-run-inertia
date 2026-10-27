package server

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
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
			s.flashError(w, "Error reading payload body")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		err = json.Unmarshal(bytes, &payload)
		if err != nil {
			s.flashError(w, "Invalid JSON Body")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		if payload.Password == "" || payload.Username == "" {
			s.flashError(w, "Both username and password are required")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		if payload.Username != s.authUser || payload.Password != s.authPassword {
			s.flashError(w, "Invalid Credentials")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		sessionKey := session.Create()

		// set a cookie on the browser!
		http.SetCookie(w, authCookie(sessionKey, 30*time.Minute))

		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}
}

// authCookie builds the server-run-auth cookie. When value is empty and
// maxAge is negative the browser is told to drop the cookie (logout).
func authCookie(value string, maxAge time.Duration) *http.Cookie {
	secure := strings.HasPrefix(os.Getenv("APP_URL"), "https://")

	c := &http.Cookie{
		Name:     "server-run-auth",
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   secure,
	}

	if maxAge < 0 {
		c.MaxAge = -1
	} else {
		c.MaxAge = int(maxAge.Seconds())
		c.Expires = time.Now().Add(maxAge)
	}

	return c
}
