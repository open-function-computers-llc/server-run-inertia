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
		createRepo(),
		addDeployKey(),
		// addWebhook(), TODO
		listRepositories(),
	}
}

// providerFlag is the shared --provider flag on the VCS commands.
func providerFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:  "provider",
		Usage: "VCS provider to use (gitea|github). Defaults to the VCS_PROVIDER env, then gitea",
	}
}
