package app

import (
	"fmt"

	"github.com/open-function-computers-llc/server-run-inertia/vcs"
	"github.com/open-function-computers-llc/server-run-inertia/vcs/gitea"
	"github.com/open-function-computers-llc/server-run-inertia/vcs/github"
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
		},
		Action: createRepoAction,
	}
}

func createRepoAction(cCtx *cli.Context) error {
	repoName := cCtx.String("name")

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
	err = p.CreateRepository(repoName)
	if err != nil {
		return err
	}
	fmt.Println("Created repo: " + repoName)
	return nil
}
