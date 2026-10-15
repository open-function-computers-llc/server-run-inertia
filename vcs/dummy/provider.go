package dummy

import "github.com/open-function-computers-llc/server-run-inertia/vcs"

type dummyProvider struct {
}

func NewProvider() (dummyProvider, error) {
	p := dummyProvider{}
	return p, nil
}

func (p dummyProvider) ListRepositories() ([]string, error) {
	return []string{}, nil
}
func (p dummyProvider) CreateRepository(name string) error {
	return vcs.ErrorNoProviderConfigured
}
