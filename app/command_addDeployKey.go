package app

import (
	"fmt"

	"github.com/open-function-computers-llc/server-run-inertia/vcs/registry"
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
				Usage: "True if this key can only READ from the repository",
			},
			providerFlag(),
		},
		Action: addDeployKeyAction,
	}
}

func addDeployKeyAction(cCtx *cli.Context) error {
	repo := cCtx.String("repository")
	name := cCtx.String("name")
	key := cCtx.String("key")
	readOnly := cCtx.Bool("read-only")

	providerName, err := registry.Resolve(cCtx.String("provider"))
	if err != nil {
		return err
	}

	p, err := registry.Build(providerName)
	if err != nil {
		return err
	}

	if err := p.AddDeployKey(repo, name, key, readOnly); err != nil {
		return err
	}
	fmt.Printf("Added deploy key '%s' to repo '%s' on %s\n", name, repo, providerName)
	return nil
}
