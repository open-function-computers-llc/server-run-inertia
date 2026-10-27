package app

import (
	"fmt"

	"github.com/open-function-computers-llc/server-run-inertia/vcs/registry"
	"github.com/urfave/cli/v2"
)

func createRepo() *cli.Command {
	return &cli.Command{
		Name:    "create-repo",
		Aliases: []string{"cr"},
		Usage:   "Create a repository in your available VCS system",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "name",
				Usage:    "Name of the new git repository",
				Required: true,
			},
			providerFlag(),
		},
		Action: createRepoAction,
	}
}

func createRepoAction(cCtx *cli.Context) error {
	repoName := cCtx.String("name")

	name, err := registry.Resolve(cCtx.String("provider"))
	if err != nil {
		return err
	}

	p, err := registry.Build(name)
	if err != nil {
		return err
	}

	if err := p.CreateRepository(repoName); err != nil {
		return err
	}
	fmt.Printf("Created repo %q on %s\n", repoName, name)
	return nil
}
