package github

func (p gitHubProvider) ListRepositories() ([]string, error) {
	return []string{"github-kurtis-repo", "github-escher-repo"}, nil
}
