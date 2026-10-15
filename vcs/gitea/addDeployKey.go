package gitea

import (
	"os"
	"strings"

	"github.com/open-function-computers-llc/server-run-inertia/vcs/gitea/client"

	gt "code.gitea.io/sdk/gitea"
)

func (p gitTeaProvider) AddDeployKey(repo, name, key string, readOnly bool) error {
	c, err := client.Client()
	if err != nil {
		return err
	}

	// repos are listed by full name (owner/repo), so honor an owner if one was given
	org := os.Getenv("GITEA_ORG")
	var owner string

	if o, r, found := strings.Cut(repo, "/"); found {
		owner = o
		repo = r
	} else if org != "" {
		owner = org
	} else {
		// for personal repos, use the authenticated user
		user, _, err := c.GetMyUserInfo()
		if err != nil {
			return err
		}
		owner = user.UserName
	}

	_, _, err = c.CreateDeployKey(owner, repo, gt.CreateKeyOption{
		Title:    name,
		Key:      key,
		ReadOnly: readOnly,
	})
	if err != nil {
		return err
	}

	return nil
}
