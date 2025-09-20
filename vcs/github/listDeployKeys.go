package github

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/open-function-computers-llc/server-run-inertia/vcs"
	"github.com/open-function-computers-llc/server-run-inertia/vcs/github/github_api"
)

func (p gitHubProvider) ListDeployKeys(repository string) ([]vcs.DeployKey, error) {

	// Relevant API docs:
	// https://docs.github.com/en/rest/deploy-keys/deploy-keys?apiVersion=2022-11-28#list-deploy-keys

	// GitHub fine-tuned token permissions required:
	// Administrator (read/write)

	if !strings.Contains(repository, "/") {
		return nil, errors.New("repository must be in the format of <owner>/<repository>")
	}

	url := "https://api.github.com/repos/" + repository + "/keys"

	// Use the client to make the request
	keys, err := github_api.GitHubMakeRequest[[]map[string]any](url, "GET", nil)
	if err != nil {
		return nil, err
	}

	// Map the keys to the DeployKey struct
	deployKeys := []vcs.DeployKey{}
	for _, key := range keys {

		createdAt, err := time.Parse(time.RFC3339, key["created_at"].(string))
		if err != nil {
			return nil, err
		}

		deployKeys = append(deployKeys, vcs.DeployKey{
			Name:      key["title"].(string),
			Key:       key["key"].(string),
			ReadOnly:  key["read_only"].(bool),
			CreatedAt: createdAt,
			CreatedBy: key["added_by"].(string),
			Id:        strconv.FormatFloat(key["id"].(float64), 'f', -1, 64),
		})
	}

	return deployKeys, nil

}
