package server

import (
	"io/fs"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/open-function-computers-llc/server-run-inertia/account"
	"github.com/open-function-computers-llc/server-run-inertia/session"
	"github.com/petaki/inertia-go"
	"github.com/sirupsen/logrus"
)

type server struct {
	port           int
	logger         *logrus.Logger
	router         *http.ServeMux
	inertiaManager *inertia.Inertia
	distFS         fs.FS
	sessions       *session.SessionBag
	authUser       string
	authPassword   string
	accounts       []account.Account
	upgrader       websocket.Upgrader
	defaultLogo    []byte
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

	// share the live reload port with the frontend
	liveReloadPort := os.Getenv("APP_LIVERELOAD_PORT")
	if liveReloadPort != "" {
		s.inertiaManager.Share("liveReloadPort", liveReloadPort)
	}

	s.bindRoutes()

	return s, nil
}

func (s *server) Serve() error {
	s.logger.Info("now serving at :" + strconv.Itoa(s.port))
	return http.ListenAndServe(":"+strconv.Itoa(s.port), s.router)
}
