package github

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/open-function-computers-llc/server-run-inertia/vcs"
	"github.com/open-function-computers-llc/server-run-inertia/vcs/github/github_api"
)

func (p gitHubProvider) AddDeployKey(repository string, key string, name string, readOnly bool) (vcs.DeployKey, error) {
	// Relevant API docs:
	// https://docs.github.com/en/rest/deploy-keys/deploy-keys?apiVersion=2022-11-28#create-a-deploy-key

	// GitHub fine-tuned token permissions required:
	// Administrator (read/write)

	if !strings.Contains(repository, "/") {
		return vcs.DeployKey{}, errors.New("repository must be in the format of <owner>/<repository>")
	}

	url := "https://api.github.com/repos/" + repository + "/keys"

	// Construct the request body
	body := map[string]any{
		"title":     name,
		"key":       key,
		"read_only": readOnly,
	}
	bodyBytes, _ := json.Marshal(body)

	// Use the client to make the request
	resp, err := github_api.GitHubMakeRequest[map[string]any](url, "POST", bodyBytes)
	if err != nil {
		return vcs.DeployKey{}, err
	}

	createdAt, err := time.Parse(time.RFC3339, resp["created_at"].(string))
	if err != nil {
		return vcs.DeployKey{}, err
	}

	return vcs.DeployKey{
		Name:      resp["title"].(string),
		Key:       resp["key"].(string),
		ReadOnly:  resp["read_only"].(bool),
		CreatedAt: createdAt,
		CreatedBy: resp["added_by"].(string),
		Id:        strconv.FormatFloat(resp["id"].(float64), 'f', -1, 64),
	}, nil
}
