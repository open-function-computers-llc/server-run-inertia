package github

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/open-function-computers-llc/server-run-inertia/vcs/github/github_api"
)

func (p gitHubProvider) CreateRepository(name string) error {

	// Relevant API docs:
	// https://docs.github.com/en/rest/repos/repos?apiVersion=2022-11-28#create-a-repository-for-the-authenticated-user
	// https://docs.github.com/en/rest/repos/repos?apiVersion=2022-11-28#create-an-organization-repository

	// GitHub fine-tuned token permissions required:
	// Metadata (read)
	// Administrator (read/write)

	url := "https://api.github.com/user/repos"
	if os.Getenv("GITHUB_ORG") != "" {
		url = "https://api.github.com/orgs/" + os.Getenv("GITHUB_ORG") + "/repos"
	}

	// Create the request body
	body := map[string]any{
		"name":        name,
		"description": "This was created automatically from server-run",
		"private":     true,
		"auto_init":   true,
	}
	bodyBytes, _ := json.Marshal(body)

	// Use the client to make the request
	repo, err := github_api.GitHubMakeRequest[map[string]any](url, "POST", bodyBytes)
	if err != nil {
		return err
	}

	// (for future reference) Other keys that might be useful on repo:
	// name, full_name, id (int), private (bool), html_url, created_at, updated_at, language

	fmt.Println("Created repo: " + repo["full_name"].(string) + " at " + repo["html_url"].(string))

	return nil
}
