package app

import (
	"fmt"

	"github.com/open-function-computers-llc/server-run-inertia/vcs"
	"github.com/open-function-computers-llc/server-run-inertia/vcs/dummy"
	"github.com/open-function-computers-llc/server-run-inertia/vcs/gitea"
	"github.com/open-function-computers-llc/server-run-inertia/vcs/github"
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
		},
		Action: addDeployKeyAction,
	}
}

func addDeployKeyAction(cCtx *cli.Context) error {
	repo := cCtx.String("repository")
	name := cCtx.String("name")
	key := cCtx.String("key")
	readOnly := cCtx.Bool("read-only")

	provider, err := vcs.DetermineProvider()
	if err != nil {
		return err
	}

	var p vcs.Provider
	switch provider {
	case "gitea":
		p, err = gitea.NewProvider()
	case "github":
		p, err = github.NewProvider()
	default:
		p, _ = dummy.NewProvider()
	}
	if err != nil {
		return err
	}

	err = p.AddDeployKey(repo, name, key, readOnly)
	if err != nil {
		return err
	}
	fmt.Println("Added deploy key '" + name + "' to repo '" + repo + "'")
	return nil
}
