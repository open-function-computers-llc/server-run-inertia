package vcs

type Provider interface {
	CreateRepository(name string) error
	ListRepositories() ([]string, error)
	AddDeployKey(repo, name, key string, readOnly bool) error
}
