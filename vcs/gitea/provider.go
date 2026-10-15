package gitea

import (
	"fmt"
	"os"

	"github.com/open-function-computers-llc/server-run-inertia/errx"
	"github.com/open-function-computers-llc/server-run-inertia/vcs"
)

type gitTeaProvider struct {
	AUTH_KEY string
	URL      string
}

func NewProvider() (gitTeaProvider, error) {
	p := gitTeaProvider{}

	for _, env := range []string{"GITEA_URL", "GITEA_TOKEN"} {
		if os.Getenv(env) == "" {
			return p, errx.Wrapf("configuring gitea", fmt.Errorf("%s is not set", env), vcs.ErrorMissingConfig)
		}
	}

	return p, nil
}
