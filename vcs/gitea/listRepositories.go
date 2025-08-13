package gitea

func (p gitTeaProvider) ListRepositories() ([]string, error) {
	return []string{"kurtis-repo", "escher-repo"}, nil
}
