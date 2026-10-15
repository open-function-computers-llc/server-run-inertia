package dummy

import "github.com/open-function-computers-llc/server-run-inertia/vcs"

func (p dummyProvider) AddDeployKey(repo, name, key string, readOnly bool) error {
	return vcs.ErrorNoProviderConfigured
}
