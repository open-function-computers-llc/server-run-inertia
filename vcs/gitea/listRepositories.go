package gitea

import (
	"os"

	"github.com/open-function-computers-llc/server-run-inertia/vcs/gitea/client"

	gt "code.gitea.io/sdk/gitea"
)

func (p gitTeaProvider) ListRepositories() ([]string, error) {
	c, err := client.Client()
	if err != nil {
		return nil, err
	}

	options := gt.ListOrgReposOptions{
		ListOptions: gt.ListOptions{
			Page: -1,
		},
	}

	output := []string{}

	if os.Getenv("GITEA_ORG") != "" {
		repos, _, err := c.ListOrgRepos(os.Getenv("GITEA_ORG"), options)
		if err != nil {
			return nil, err
		}
		for _, r := range repos {
			output = append(output, r.FullName)
		}
	} else {
		// TODO: list personal repos
	}
	return output, nil
}
