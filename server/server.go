package server

import (
	"io/fs"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/open-function-computers-llc/server-run-inertia/account"
	"github.com/open-function-computers-llc/server-run-inertia/session"
	"github.com/open-function-computers-llc/server-run-inertia/vcs"
	"github.com/open-function-computers-llc/server-run-inertia/vcs/dummy"
	"github.com/open-function-computers-llc/server-run-inertia/vcs/registry"
	"github.com/petaki/inertia-go"
	"github.com/sirupsen/logrus"
)

type server struct {
	port            int
	logger          *logrus.Logger
	router          *http.ServeMux
	inertiaManager  *inertia.Inertia
	distFS          fs.FS
	sessions        *session.SessionBag
	authUser        string
	authPassword    string
	accounts        []account.Account
	upgrader        websocket.Upgrader
	defaultLogo     []byte
	vcsProvider     vcs.Provider
	vcsProviders    []string
	vcsErrorMessage string
	vcsAvailable    bool
	cms             string
	socketPath      string
}

func New(port int, url string, inertiaFS fs.FS, logo []byte) (server, error) {
	s := server{
		logger:         logrus.New(),
		port:           port,
		router:         http.NewServeMux(),
		inertiaManager: inertia.NewWithFS(url, "dist/index.html", "", inertiaFS),
		distFS:         inertiaFS,
		sessions:       session.Initialize(),
		authUser:       os.Getenv("AUTH_USER"),
		authPassword:   os.Getenv("AUTH_PASSWORD"),
		defaultLogo:    logo,
	}

	// set debug loggin if we are in dev mode
	if os.Getenv("APP_ENV") == "development" {
		s.logger.SetLevel(logrus.DebugLevel)
	}

	err := s.bootstrapAccounts()
	if err != nil {
		return s, err
	}

	// share the app environment with the frontend
	appEnvironment := os.Getenv("APP_ENV")
	if appEnvironment == "" {
		appEnvironment = "production"
	}
	s.inertiaManager.Share("appEnvironment", appEnvironment)

	// share which CMS this server hosts (wordpress|modx, or "" if unknown)
	s.cms = detectCMS(url)
	s.inertiaManager.Share("cms", s.cms)

	// share the live reload port with the frontend
	liveReloadPort := os.Getenv("APP_LIVERELOAD_PORT")
	if liveReloadPort != "" {
		s.inertiaManager.Share("liveReloadPort", liveReloadPort)
	}

	s.vcsProvider, _ = dummy.NewProvider()

	// check which VCS providers are configured on this instance. Providers are
	// resolved per client (account), with gitea as the default. A misconfigured
	// provider only disables VCS features, it doesn't stop the server.
	s.vcsProviders = registry.Configured()
	s.vcsAvailable = len(s.vcsProviders) > 0

	if !s.vcsAvailable {
		s.vcsErrorMessage = "Version control isn't configured on this server. Set GITEA_URL and GITEA_TOKEN (default provider) or GITHUB_TOKEN, and optionally VCS_PROVIDER to change the default."
	} else if name, err := registry.Resolve(""); err != nil {
		s.vcsErrorMessage = err.Error()
	} else if p, err := registry.Build(name); err != nil {
		s.vcsErrorMessage = err.Error()
	} else {
		s.vcsProvider = p
	}

	s.bindRoutes()

	go s.setupSocket()

	return s, nil
}

func (s *server) Serve() error {
	s.logger.Info("now serving at :" + strconv.Itoa(s.port))
	return http.ListenAndServe(":"+strconv.Itoa(s.port), s.router)
}
