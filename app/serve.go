package app

import (
	"io/fs"

	"github.com/open-function-computers-llc/server-run-inertia/server"
	"github.com/urfave/cli/v2"
)

func serve(port int, url string, dist fs.FS, logo []byte) *cli.Command {
	return &cli.Command{
		Name:    "serve",
		Aliases: []string{"s"},
		Usage:   "start up the UI web server",
		Action: func(cCtx *cli.Context) error {
			s, err := server.New(port, url, dist, logo)
			if err != nil {
				panic("Trouble setting up server: " + err.Error())
			}

			return s.Serve()
		},
	}
}
