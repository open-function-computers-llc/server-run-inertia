package app

import (
	"fmt"

	"github.com/open-function-computers-llc/server-run-inertia/vcs/registry"
	"github.com/urfave/cli/v2"
)

func listRepositories() *cli.Command {
	return &cli.Command{
		Name:    "list-repos",
		Aliases: []string{"lr"},
		Usage:   "List repositories in your available VCS system",
		Flags: []cli.Flag{
			providerFlag(),
		},
		Action: listRepoAction,
	}
}

func listRepoAction(cCtx *cli.Context) error {
	name, err := registry.Resolve(cCtx.String("provider"))
	if err != nil {
		return err
	}

	p, err := registry.Build(name)
	if err != nil {
		return err
	}

	repos, err := p.ListRepositories()
	if err != nil {
		return err
	}
	for _, repoName := range repos {
		fmt.Println(repoName)
	}

	return nil
}
