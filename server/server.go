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
	"github.com/open-function-computers-llc/server-run-inertia/vcs/gitea"
	"github.com/open-function-computers-llc/server-run-inertia/vcs/github"
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

	// check if VCS is enabled on this server run instance. A misconfigured
	// provider only disables VCS features, it doesn't stop the server.
	provider, _ := vcs.DetermineProvider()
	switch provider {
	case "gitea":
		if p, err := gitea.NewProvider(); err != nil {
			s.vcsErrorMessage = err.Error()
		} else {
			s.vcsProvider = p
			s.vcsAvailable = true
		}
	case "github":
		if p, err := github.NewProvider(); err != nil {
			s.vcsErrorMessage = err.Error()
		} else {
			s.vcsProvider = p
			s.vcsAvailable = true
		}
	case "bitbucket":
		s.vcsErrorMessage = "Bitbucket isn't supported yet"
	}

	s.bindRoutes()

	go s.setupSocket()

	return s, nil
}

func (s *server) Serve() error {
	s.logger.Info("now serving at :" + strconv.Itoa(s.port))
	return http.ListenAndServe(":"+strconv.Itoa(s.port), s.router)
}
