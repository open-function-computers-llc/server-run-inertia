package app

import (
	"fmt"

	"github.com/open-function-computers-llc/server-run-inertia/vcs"
	"github.com/open-function-computers-llc/server-run-inertia/vcs/gitea"
	"github.com/open-function-computers-llc/server-run-inertia/vcs/github"
	"github.com/urfave/cli/v2"
)

func listDeployKeys() *cli.Command {
	return &cli.Command{
		Name:    "list-deploy-keys",
		Aliases: []string{"ldk"},
		Usage:   "List deploy keys from a given repository in your available VCS system",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "repo",
				Aliases:  []string{"r"},
				Usage:    "Name of repository that we are listing deploy keys for",
				Required: true,
			},
		},
		Action: listDeployKeysAction,
	}
}

func listDeployKeysAction(cCtx *cli.Context) error {
	repository := cCtx.String("repo")

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

	deployKeys, err := p.ListDeployKeys(repository)
	if err != nil {
		return err
	}
	for _, deployKey := range deployKeys {
		fmt.Println(deployKey)
	}

	return nil
}
