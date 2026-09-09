package scriptrunner

import (
	"errors"
	"fmt"
	"os"

	"github.com/open-function-computers-llc/server-run-inertia/errx"
)

var ARGSCRIPT string = "args"
var PLAINSCRIPT string = "plain"
var ENVSCRIPT string = "env"

func scriptType(name string) (string, error) {
	for v := range argumentScripts {
		if v == name {
			return ARGSCRIPT, nil
		}
	}
	for v := range plainScripts {
		if v == name {
			return PLAINSCRIPT, nil
		}
	}
	for v := range envScripts {
		if v == name {
			return ENVSCRIPT, nil
		}
	}

	return "", errx.Wrapf("invalid script", fmt.Errorf("script %q is not a registered script", name), ErrorInvalidScript)
}

func VerifyAllScriptsExist() error {
	for _, v := range argumentScripts {
		if _, err := os.Stat(os.Getenv("SCRIPTS_ROOT") + v); errors.Is(err, os.ErrNotExist) {
			return errx.Wrapf("missing script", fmt.Errorf("script %q not found in SCRIPTS_ROOT", v), ErrorMissingScript)
		}
	}
	for _, v := range plainScripts {
		if _, err := os.Stat(os.Getenv("SCRIPTS_ROOT") + v); errors.Is(err, os.ErrNotExist) {
			return errx.Wrapf("missing script", fmt.Errorf("script %q not found in SCRIPTS_ROOT", v), ErrorMissingScript)
		}
	}
	for _, v := range envScripts {
		if _, err := os.Stat(os.Getenv("SCRIPTS_ROOT") + v.filename); errors.Is(err, os.ErrNotExist) {
			return errx.Wrapf("missing script", fmt.Errorf("script %q not found in SCRIPTS_ROOT", v.filename), ErrorMissingScript)
		}
	}
	return nil
}
