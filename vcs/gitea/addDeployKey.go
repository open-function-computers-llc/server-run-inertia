package gitea

import (
	"errors"

	"github.com/open-function-computers-llc/server-run-inertia/vcs"
)

func (p gitTeaProvider) AddDeployKey(repository string, key string, name string, readOnly bool) (vcs.DeployKey, error) {

	// TODO: Implement
	return vcs.DeployKey{}, errors.New("not implemented")
}
