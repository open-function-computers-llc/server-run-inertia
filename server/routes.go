package server

import (
	"net/http"
)

func (s *server) bindRoutes() {
	openRoutes := map[string]http.HandlerFunc{
		"/login": s.handleLogin(),
		"/":      s.handleIndex(),
	}

	for path, handler := range openRoutes {
		s.router.Handle(path, s.inertiaManager.Middleware(s.LogRequest(handler)))
	}

	protectedRoutes := map[string]http.HandlerFunc{
		"/home": s.handleHome(),
	}

	for path, handler := range protectedRoutes {
		s.router.Handle(path, s.inertiaManager.Middleware(s.LogRequest(handler)))
	}

	s.router.Handle("/dist/", http.FileServer(http.FS(s.distFS)))
}
