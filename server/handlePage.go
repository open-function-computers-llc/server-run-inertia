package server

import (
	"net/http"
)

func (s *server) handlePage(pageName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.render(w, r, pageName, nil)
	}
}
