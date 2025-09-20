package dummy

import "github.com/open-function-computers-llc/server-run-inertia/vcs"

type dummyProvider struct {
}

func NewProvider() (dummyProvider, error) {
	p := dummyProvider{}
	return p, nil
}

func (p dummyProvider) ListRepositories() ([]string, error) {
	return []string{""}, nil
}
func (p dummyProvider) CreateRepository(name string) error {
	return nil
}
func (p dummyProvider) ListDeployKeys(repository string) ([]vcs.DeployKey, error) {
	return []vcs.DeployKey{}, nil
}
func (p dummyProvider) AddDeployKey(repository string, key string, name string, readOnly bool) (vcs.DeployKey, error) {
	return vcs.DeployKey{}, nil
}
func (p dummyProvider) RemoveDeployKey(repository string, id string) error {
	return nil
}
