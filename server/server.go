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
}

func New(port int, url string, fs fs.FS) (server, error) {
	s := server{
		logger:         logrus.New(),
		port:           port,
		router:         http.NewServeMux(),
		inertiaManager: inertia.New(url, "./server/index.html", ""),
		distFS:         fs,
		sessions:       session.Initialize(),
		authUser:       os.Getenv("AUTH_USER"),
		authPassword:   os.Getenv("AUTH_PASSWORD"),
	}
	err := s.bootstrapAccounts()
	if err != nil {
		return s, err
	}

	s.bindRoutes()

	return s, nil
}

func (s *server) Serve() error {
	s.logger.Info("now serving at :" + strconv.Itoa(s.port))
	return http.ListenAndServe(":"+strconv.Itoa(s.port), s.router)
}
