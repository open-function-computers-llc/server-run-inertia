package server

import (
	"fmt"

	"github.com/open-function-computers-llc/server-run-inertia/account"
	"github.com/open-function-computers-llc/server-run-inertia/errx"
)

func (s *server) findAccount(name string) (*account.Account, error) {
	for _, a := range s.accounts {
		if a.Name == name {
			return &a, nil
		}
	}

	return nil, errx.Wrapf("account not found", fmt.Errorf("account %q is not registered with this server", name), ErrorAccountNotFound)
}
