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
		// check for other required ENV VARS for a BITBUCKET connection
		return "bitbucket", nil
	}

	if configuredProvider == "GITHUB" {
		// check for other required ENV VARS for a GITHUB connection

		if os.Getenv("GITHUB_TOKEN") == "" {
			return "", errors.New("GITHUB_TOKEN is not set. Required for GitHub VCS provider.")
		}

		return "github", nil
	}

	return "", errors.New("Invalid ENV value for VCS_PROVIDER")
}
