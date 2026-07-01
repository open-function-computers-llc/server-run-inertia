package server

import (
	"net/http"
	"strings"
)

func (s *server) bindRoutes() {
	s.logger.Debug("---")
	s.logger.Debug("Registering routes")
	s.logger.Debug("---")
	s.logger.Debug("")
	s.logger.Debug("OPEN ROUTES")
	openRoutes := map[string]http.HandlerFunc{
		"GET /login":        s.handlePage("Login"),
		"POST /handle-auth": s.handleFormProcessAuth(),
		"GET /error":        s.handlePage("Error"),
		"GET /logo":         s.handleLogo(),
		"GET /":             s.handlePage("Index"),
	}

	for path, handler := range openRoutes {
		prefix := "Registering open route"
		pathParts := strings.Split(path, " ")
		if len(pathParts) < 2 {
			continue
		}

		s.logger.Debug(alignRouteLog(prefix, pathParts[0]), pathParts[1])
		s.router.Handle(path, s.inertiaManager.Middleware(s.LogRequest(handler)))
	}
	s.logger.Debug("")
	s.logger.Debug("PROTECTED ROUTES")

	protectedRoutes := map[string]http.HandlerFunc{
		// inertia pages
		"GET /dashboard":                              s.handleDashboard(),
		"GET /create-account":                         s.handlePage("Account/Create"),
		"GET /account/{account}":                      s.handleAccountDetails(""),
		"GET /importable-accounts":                    s.handleListImportableAccounts(),
		"GET /logout":                                 s.handleLogout(),
		"GET /settings":                               s.handleSettings(),
		"GET /vcs":                                    s.handleVcs(),
		"POST /accounts/{account}/domains":            s.handleAddDomain(),
		"PUT /accounts/{account}/domains/primary":     s.handleSetPrimaryDomain(),
		"DELETE /accounts/{account}/domains/{domain}": s.handleDeleteDomain(),
		"PATCH /accounts/{account}/uptime-url":        s.handleFormSetUptimeURL(),

		// raw JSON routes
		"GET /api/accounts/{account}/analytics": s.handleAccountAnalytics(),
		"POST /api/get-uptime-info":             s.handleUptimeInfo(),
	}

	for path, handler := range protectedRoutes {
		prefix := "Registering protected route"
		pathParts := strings.Split(path, " ")
		if len(pathParts) < 2 {
			continue
		}

		s.logger.Debug(alignRouteLog(prefix, pathParts[0]), pathParts[1])
		s.router.Handle(path, s.inertiaManager.Middleware(s.LogRequest(s.ProtectRequest(handler))))
	}
	s.logger.Debug("")
	s.logger.Debug("EXTERNAL ROUTES")

	// external api access for report generator or other 3rd party access
	// all routes are GET only
	thirdPartyAccess := map[string]http.HandlerFunc{
		"GET /analytics": s.handleThirdPartyAnalyticsJSON(),
	}
	for path, handler := range thirdPartyAccess {
		prefix := "Registering external route"
		pathParts := strings.Split(path, " ")
		if len(pathParts) < 2 {
			continue
		}

		// add the route prefix to these special cases
		pathParts[1] = "/external-access" + pathParts[1]

		s.logger.Debug(alignRouteLog(prefix, pathParts[0]), pathParts[1])
		s.router.Handle(pathParts[0]+" "+pathParts[1], s.LogRequest(s.thirdPartyProtection(handler)))
	}
	s.logger.Debug("")
	s.logger.Debug("SOCKET ROUTES")

	// websocket routes
	websocketRoutes := map[string]http.HandlerFunc{
		"GET /system-load":   s.handleAPISystemLoad(),
		"GET /script-runner": s.handleStreamScriptRunner(),
	}
	for path, handler := range websocketRoutes {
		prefix := "Registering web socket route"
		pathParts := strings.Split(path, " ")
		if len(pathParts) < 2 {
			continue
		}

		// add the route prefix to these special cases
		pathParts[1] = "/stream" + pathParts[1]

		s.logger.Debug(alignRouteLog(prefix, pathParts[0]), pathParts[1])
		s.router.Handle(pathParts[0]+" "+pathParts[1], s.LogRequest(s.ProtectRequest(handler)))
	}
	s.logger.Debug("")
	s.logger.Debug("MOUNTING ROUTE FILESYSTEM AT: /dist/*")

	s.router.Handle("GET /dist/", http.FileServer(http.FS(s.distFS)))

	s.logger.Debug("---")
	s.logger.Debug("Route binding complete!")
	s.logger.Debug("---")
}

func alignRouteLog(prefix, verb string) string {
	extraAlignmentSpaces := ""
	if verb == "PATCH" {
		extraAlignmentSpaces = " "
	}
	if verb == "POST" {
		extraAlignmentSpaces = "  "
	}
	if verb == "GET" || verb == "PUT" {
		extraAlignmentSpaces = "   "
	}
	return prefix + ": (" + verb + ") " + extraAlignmentSpaces
}
