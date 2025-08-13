package app

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func addDeployKey() *cli.Command {
	return &cli.Command{
		Name:    "add-deploy-key",
		Aliases: []string{"adk"},
		Usage:   "Add a deploy key to a given git repository",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "repository",
				Usage:    "Name of repository that we are adding a deploy key to",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "name",
				Usage:    "Name of the new deploy key",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "key",
				Usage:    "Public RSA key",
				Required: true,
			},
			&cli.BoolFlag{
				Name:  "read-only",
				Usage: "False if this key can PUSH to the repository",
			},
		},
		Action: func(cCtx *cli.Context) error {
			fmt.Println("TODO!")
			return nil
		},
	}
}
