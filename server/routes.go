package server

import (
	"net/http"
)

func (s *server) bindRoutes() {
	openRoutes := map[string]http.HandlerFunc{
		"GET /login":        s.handlePage("Login"),
		"POST /handle-auth": s.handleProcessAuth(),
		"GET /error":        s.handlePage("Error"),
		"/":                 s.handlePage("Index"),
	}

	for path, handler := range openRoutes {
		s.router.Handle(path, s.inertiaManager.Middleware(s.LogRequest(handler)))
	}

	protectedRoutes := map[string]http.HandlerFunc{
		// inertia pages
		"GET /dashboard":             s.handleDashboard(),
		"GET /accounts":              s.handleAccounts(),
		"GET /create-account":        s.handlePage("Account/Create"),
		"GET /account/{accountName}": s.handleAccountDetails(),
	}

	for path, handler := range protectedRoutes {
		s.router.Handle(path, s.inertiaManager.Middleware(s.LogRequest(s.ProtectRequest(handler))))
	}

	// external api access for report generator or other 3rd party access
	// all routes are GET only
	thirdPartyAccess := map[string]http.HandlerFunc{
		"analytics": s.handleThirdPartyAnalyticsJSON(),
	}
	for path, handler := range thirdPartyAccess {
		s.router.Handle("GET /external-access/"+path, s.LogRequest(s.thirdPartyProtection(handler)))
	}

	// websocket routes
	websocketRoutes := map[string]http.HandlerFunc{
		"GET /stream/system-load":   s.handleAPISystemLoad(),
		"GET /stream/script-runner": s.handleStreamScriptRunner(),
	}
	for path, handler := range websocketRoutes {
		s.router.Handle(path, s.ProtectRequest(s.LogRequest(handler)))
	}

	s.router.Handle("/dist/", http.FileServer(http.FS(s.distFS)))
}
