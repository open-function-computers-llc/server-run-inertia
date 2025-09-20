package app

import (
	"fmt"

	"github.com/open-function-computers-llc/server-run-inertia/vcs"
	"github.com/open-function-computers-llc/server-run-inertia/vcs/gitea"
	"github.com/open-function-computers-llc/server-run-inertia/vcs/github"
	"github.com/urfave/cli/v2"
)

func removeDeployKey() *cli.Command {
	return &cli.Command{
		Name:    "remove-deploy-key",
		Aliases: []string{"rdk"},
		Usage:   "Remove a deploy key from a given git repository",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "repo",
				Aliases:  []string{"r"},
				Usage:    "Name of repository that we are removing a deploy key from",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "id",
				Usage:    "ID of the deploy key to remove",
				Required: true,
			},
		},
		Action: func(cCtx *cli.Context) error {
			repository := cCtx.String("repo")
			id := cCtx.String("id")

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

			err = p.RemoveDeployKey(repository, id)
			if err != nil {
				return err
			}
			fmt.Println("Removed deploy key from " + repository + " with id " + id)
			return nil
		},
	}
}
