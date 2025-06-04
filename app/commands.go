package app

import (
	"io/fs"

	"github.com/urfave/cli/v2"
)

func AvailableCommands(port int, url string, dist fs.FS, logo []byte) []*cli.Command {
	return []*cli.Command{
		serve(port, url, dist, logo),
		listSetting(),
		envExport(),
		update(),
	}
}
