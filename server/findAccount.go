package server

import (
	"errors"

	"github.com/open-function-computers-llc/server-run-inertia/account"
)

func (s *server) findAccount(name string) (*account.Account, error) {
	for _, a := range s.accounts {
		if a.Name == name {
			return &a, nil
		}
	}

	return nil, errors.New("couldn't find the account " + name)
}
