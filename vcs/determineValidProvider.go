package vcs

import (
	"errors"
	"os"
)

func DetermineProvider() (string, error) {
	configuredProvider := os.Getenv("VCS_PROVIDER")

	if configuredProvider == "" {
		return "", errors.New("No VCS provider found in ENV. Make sure VCS_PROVIDER is set.")
	}

	if configuredProvider == "GITEA" {
		// check for other required ENV VARS for a GITEA connection
		return "gitea", nil
	}

	if configuredProvider == "BITBUCKET" {
		// check for other required ENV VARS for a GITEA connection
		return "bitbucket", nil
	}

	if configuredProvider == "GITHUB" {
		// check for other required ENV VARS for a GITEA connection
		return "github", nil
	}

	return "", errors.New("Invalid ENV value for VCS_PROVIDER")
}
