package app

import (
	"fmt"

	"github.com/open-function-computers-llc/server-run-inertia/vcs"
	"github.com/open-function-computers-llc/server-run-inertia/vcs/gitea"
	"github.com/open-function-computers-llc/server-run-inertia/vcs/github"
	"github.com/urfave/cli/v2"
)

func listRepositories() *cli.Command {
	return &cli.Command{
		Name:    "list-repos",
		Aliases: []string{"lr"},
		Usage:   "List repositories in your available VCS system",
		Action:  listRepoAction,
	}
}

func listRepoAction(cCtx *cli.Context) error {
	provider, err := vcs.DetermineProvider()
	if err != nil {
		return err
	}

	var p vcs.Provider
	if provider == "gitea" {
		p, err = gitea.NewProvider()
		if err != nil {
			return err
		}
	}
	if provider == "github" {
		p, err = github.NewProvider()
		if err != nil {
			return err
		}
	}
	if provider == "bitbucket" {
		// p = bitbucket.NewProvider()
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
