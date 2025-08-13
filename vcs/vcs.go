package vcs

type Provider interface {
	CreateRepository(name string) error
	ListRepositories() ([]string, error)
}
