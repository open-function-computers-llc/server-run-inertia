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

	org := os.Getenv("GITEA_ORG")
	output := []string{}

	// follow the server's pagination rather than assuming a page size
	for page := 1; page > 0; {
		var repos []*gt.Repository
		var resp *gt.Response

		listOptions := gt.ListOptions{Page: page, PageSize: 50}
		if org != "" {
			repos, resp, err = c.ListOrgRepos(org, gt.ListOrgReposOptions{ListOptions: listOptions})
		} else {
			repos, resp, err = c.ListMyRepos(gt.ListReposOptions{ListOptions: listOptions})
		}
		if err != nil {
			return nil, err
		}

		for _, r := range repos {
			output = append(output, r.FullName)
		}
		page = resp.NextPage
	}

	return output, nil
}
