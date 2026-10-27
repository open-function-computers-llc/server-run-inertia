// Package registry builds VCS providers. A provider is chosen per client
// (account), not per server: the account's own setting wins, then the
// VCS_PROVIDER env, then the default (gitea) if it's configured.
package registry

import (
	"fmt"
	"slices"

	"github.com/open-function-computers-llc/server-run-inertia/vcs"
	"github.com/open-function-computers-llc/server-run-inertia/vcs/dummy"
	"github.com/open-function-computers-llc/server-run-inertia/vcs/gitea"
	"github.com/open-function-computers-llc/server-run-inertia/vcs/github"
)

// DefaultProvider is used when an account hasn't picked one and the server
// hasn't set VCS_PROVIDER.
const DefaultProvider = "gitea"

// Settable are the providers an account can choose.
var Settable = []string{"gitea", "github"}

// Build constructs a single provider by name.
func Build(name string) (vcs.Provider, error) {
	switch name {
	case "gitea":
		p, err := gitea.NewProvider()
		if err != nil {
			return nil, err
		}
		return p, nil
	case "github":
		p, err := github.NewProvider()
		if err != nil {
			return nil, err
		}
		return p, nil
	case "bitbucket":
		return nil, fmt.Errorf("bitbucket isn't supported yet: %w", vcs.ErrorNotImplemented)
	case "":
		p, _ := dummy.NewProvider()
		return p, nil
	}

	return nil, fmt.Errorf("unknown VCS provider %q: %w", name, vcs.ErrorMissingConfig)
}

// Configured returns every provider whose required env vars are present.
func Configured() []string {
	configured := []string{}

	if _, err := gitea.NewProvider(); err == nil {
		configured = append(configured, "gitea")
	}
	if _, err := github.NewProvider(); err == nil {
		configured = append(configured, "github")
	}

	return configured
}

// Resolve works out which provider name should serve a request.
// accountProvider is the account's own setting ("" when unset).
func Resolve(accountProvider string) (string, error) {
	if accountProvider != "" {
		if !slices.Contains(Settable, accountProvider) {
			return "", fmt.Errorf("invalid VCS provider %q for this account", accountProvider)
		}
		return accountProvider, nil
	}

	if fromEnv, err := vcs.DetermineProvider(); err == nil {
		return fromEnv, nil
	}

	configured := Configured()
	if slices.Contains(configured, DefaultProvider) {
		return DefaultProvider, nil
	}
	if len(configured) > 0 {
		return configured[0], nil
	}

	return "", vcs.ErrorNoProviderConfigured
}

// ForAccount resolves and builds the provider that should serve an account,
// returning the provider along with the name it was built from.
func ForAccount(accountProvider string) (vcs.Provider, string, error) {
	name, err := Resolve(accountProvider)
	if err != nil {
		return nil, "", err
	}

	p, err := Build(name)
	if err != nil {
		return nil, "", err
	}

	return p, name, nil
}
