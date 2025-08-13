package vcs

import (
	"os"
	"testing"
)

func Test_WeCanDetermineTheCorrectVcsProviderFromTheCurrentAppEnv(t *testing.T) {
	_, err := DetermineProvider()
	if err == nil {
		t.Error("No ENV means there shouldn't be a provider!")
	}

	os.Setenv("VCS_PROVIDER", "thisisntreal") // not a valid provider
	_, err = DetermineProvider()
	if err == nil {
		t.Error("You must have a valid option in your env. GITEA|GITHUB|BITBUCKET")
	}

	os.Setenv("VCS_PROVIDER", "GITEA")
	provider, err := DetermineProvider()
	if provider != "gitea" {
		t.Error("Provider should have been set!")
	}

	os.Setenv("VCS_PROVIDER", "GITHUB")
	provider, err = DetermineProvider()
	if provider != "github" {
		t.Error("Provider should have been set!")
	}

	os.Setenv("VCS_PROVIDER", "BITBUCKET")
	provider, err = DetermineProvider()
	if provider != "bitbucket" {
		t.Error("Provider should have been set!")
	}
}
