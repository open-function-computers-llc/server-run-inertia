package github

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
)

func (p gitHubProvider) ListRepositories() ([]string, error) {

	url := "https://api.github.com/user/repos"

	// If an org is set, use the org's repos
	if os.Getenv("GITHUB_ORG") != "" {
		url = "https://api.github.com/orgs/" + os.Getenv("GITHUB_ORG") + "/repos"
	}

	// Create the request and set the authorization header
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "token "+os.Getenv("GITHUB_TOKEN"))
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Check if the request was successful
	if resp.StatusCode != 200 {
		// Parse the error response
		var errorResponse map[string]interface{}
		err = json.Unmarshal(body, &errorResponse)
		if err == nil && errorResponse["message"] != nil {
			return nil, errors.New("GitHub API Error: " + errorResponse["message"].(string))
		}
		// Generic error message
		return nil, errors.New("Failed to list repositories")
	}

	// Parse the response body
	var repos []map[string]interface{}
	err = json.Unmarshal(body, &repos)
	if err != nil {
		return nil, err
	}

	// Extract just the repository names
	repoNames := []string{}
	for _, repo := range repos {
		repoNames = append(repoNames, repo["full_name"].(string))
		// (future reference) Other keys that might be useful:
		// id (int), private (bool), url, created_at, updated_at, language
	}

	return repoNames, nil
}
