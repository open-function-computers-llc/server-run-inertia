package github

import (
	"errors"
	"strings"

	"github.com/open-function-computers-llc/server-run-inertia/vcs/github/github_api"
)

func (p gitHubProvider) RemoveDeployKey(repository string, id string) error {

	// Relevant API docs:
	// https://docs.github.com/en/rest/deploy-keys/deploy-keys?apiVersion=2022-11-28#delete-a-deploy-key

	// GitHub fine-tuned token permissions required:
	// Administrator (read/write)

	if !strings.Contains(repository, "/") {
		return errors.New("repository must be in the format of <owner>/<repository>")
	}

	url := "https://api.github.com/repos/" + repository + "/keys/" + id

	// Use the client to make the request
	_, err := github_api.GitHubMakeRequest[map[string]any](url, "DELETE", nil)
	if err != nil {
		return err
	}

	return nil
}
