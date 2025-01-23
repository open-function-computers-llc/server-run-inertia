package scriptrunner

import (
	"errors"
	"os"
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

	return "", errors.New(InvalidScript + " called: " + name)
}

func VerifyAllScriptsExist() error {
	for _, v := range argumentScripts {
		if _, err := os.Stat(os.Getenv("SCRIPTS_ROOT") + v); errors.Is(err, os.ErrNotExist) {
			return errors.New(MissingScript + ": " + v)
		}
	}
	for _, v := range plainScripts {
		if _, err := os.Stat(os.Getenv("SCRIPTS_ROOT") + v); errors.Is(err, os.ErrNotExist) {
			return errors.New(MissingScript + ": " + v)
		}
	}
	for _, v := range envScripts {
		if _, err := os.Stat(os.Getenv("SCRIPTS_ROOT") + v); errors.Is(err, os.ErrNotExist) {
			return errors.New(MissingScript + ": " + v)
		}
	}
	return nil
}
