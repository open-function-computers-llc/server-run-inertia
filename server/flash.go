package server

import (
	"encoding/base64"
	"net/http"
)

const flashErrorCookie = "server-run-flash-error"

// flashError stores a one-time error message in a short-lived cookie so it is
// shown on the next rendered page for this browser only. Use it right before
// redirecting.
func (s *server) flashError(w http.ResponseWriter, message string) {
	http.SetCookie(w, &http.Cookie{
		Name:     flashErrorCookie,
		Value:    base64.RawURLEncoding.EncodeToString([]byte(message)),
		Path:     "/",
		MaxAge:   60,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

// render renders an inertia page, adding (and then clearing) any flashed error
// as the "error" prop.
func (s *server) render(w http.ResponseWriter, r *http.Request, component string, props map[string]any) error {
	if props == nil {
		props = map[string]any{}
	}

	if message, ok := s.consumeFlashError(w, r); ok {
		if _, exists := props["error"]; !exists {
			props["error"] = message
		}
	}

	return s.inertiaManager.Render(w, r, component, props)
}

func (s *server) consumeFlashError(w http.ResponseWriter, r *http.Request) (string, bool) {
	c, err := r.Cookie(flashErrorCookie)
	if err != nil || c.Value == "" {
		return "", false
	}

	http.SetCookie(w, &http.Cookie{
		Name:     flashErrorCookie,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	message, err := base64.RawURLEncoding.DecodeString(c.Value)
	if err != nil {
		return "", false
	}
	return string(message), true
}
