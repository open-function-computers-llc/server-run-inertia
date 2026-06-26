package server

import (
	"net/http"
	"strconv"
	"time"

	"github.com/open-function-computers-llc/server-run-inertia/session"
)

func (s *server) ProtectRequest(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie("server-run-auth")
		if err != nil {
			ctx := s.inertiaManager.WithProp(r.Context(), "errMessage", "no cookie")
			http.Redirect(w, r.WithContext(ctx), "/login", http.StatusSeeOther)
			return
		}

		if sessionCookie.Value == "" {
			ctx := s.inertiaManager.WithProp(r.Context(), "errMessage", "no cookie")
			http.Redirect(w, r.WithContext(ctx), "/login", http.StatusSeeOther)
			return
		}

		cookieIsValid, validUntil := session.Validate(s.logger, sessionCookie.Value)
		if !cookieIsValid {
			ctx := s.inertiaManager.WithProp(r.Context(), "errMessage", "invalid cookie")
			http.Redirect(w, r.WithContext(ctx), "/login", http.StatusSeeOther)
			return
		}

		ctx := s.inertiaManager.WithProp(r.Context(), "expiresIn", strconv.Itoa(int(time.Until(validUntil).Seconds())))
		next(w, r.WithContext(ctx))
	}
}
