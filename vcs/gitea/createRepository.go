package gitea

import (
	"os"

	"github.com/open-function-computers-llc/server-run-inertia/vcs/gitea/client"

	gt "code.gitea.io/sdk/gitea"
)

func (p gitTeaProvider) CreateRepository(name string) error {
	c, err := client.Client()
	if err != nil {
		return err
	}

	options := gt.CreateRepoOption{
		Name:          name,
		Description:   "This was created automatically from server-run",
		Private:       true,
		AutoInit:      true,
		DefaultBranch: "main",
	}

	if os.Getenv("GITEA_ORG") != "" {
		_, _, err = c.CreateOrgRepo(os.Getenv("GITEA_ORG"), options)
		if err != nil {
			return err
		}
	} else {
		_, _, err = c.CreateRepo(options)
		if err != nil {
			return err
		}
	}

	return nil
}
