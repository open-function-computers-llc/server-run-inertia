package server

import (
	"os"
	"os/exec"
	"strings"

	"github.com/open-function-computers-llc/server-run-inertia/account"
)

func (s *server) bootstrapAccounts() error {
	s.accounts = []account.Account{}

	stdout, err := exec.Command("ls", os.Getenv("WEBSITES_ROOT")).Output()
	if err != nil {
		return err
	}
	lines := strings.Split(string(stdout), "\n")
	for _, line := range lines {
		if len(line) == 0 { // skip blank lines coming from the `ls` command
			continue
		}
		a, err := account.New(line)
		if err != nil {
			return err
		}
		err = a.LoadStatus()
		if err != nil {
			return err
		}

		s.accounts = append(s.accounts, a)
	}

	return nil
}
