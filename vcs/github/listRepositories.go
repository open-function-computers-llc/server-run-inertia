package github

import (
	"os"

	"github.com/open-function-computers-llc/server-run-inertia/vcs/github/github_api"
)

func (p gitHubProvider) ListRepositories() ([]string, error) {

	// Relevant API docs:
	// https://docs.github.com/en/rest/repos/repos?apiVersion=2022-11-28#list-repositories-for-the-authenticated-user
	// https://docs.github.com/en/rest/repos/repos?apiVersion=2022-11-28#list-organization-repositories

	// GitHub fine-tuned token permissions required:
	// Metadata (read)

	url := "https://api.github.com/user/repos"

	// If an org is set, use the org's repos
	if os.Getenv("GITHUB_ORG") != "" {
		url = "https://api.github.com/orgs/" + os.Getenv("GITHUB_ORG") + "/repos"
	}

	// Use the client to make the request
	repos, err := github_api.GitHubMakeRequest[[]map[string]any](url, "GET", nil)
	if err != nil {
		return nil, err
	}

	// Extract just the repository names
	repoNames := []string{}
	for _, repo := range repos {
		repoNames = append(repoNames, repo["full_name"].(string))
		// (future reference) Other keys that might be useful:
		// id (int), private (bool), html_url, url, created_at, updated_at, language
	}

	return repoNames, nil
}
