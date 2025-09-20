package gitea

import "github.com/open-function-computers-llc/server-run-inertia/vcs"

func (p gitTeaProvider) ListDeployKeys(repository string) ([]vcs.DeployKey, error) {

	// TODO: Implement
	return []vcs.DeployKey{}, nil
}
