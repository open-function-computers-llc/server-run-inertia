package scriptrunner

import (
	"fmt"
	"net"
	"regexp"

	"github.com/open-function-computers-llc/server-run-inertia/errx"
)

var envKeyPattern = regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)

// validateEnv checks a script's incoming env vars against its registration:
// required vars are present, unregistered vars are rejected, keys look like
// env vars, and any per-script format checks pass.
func validateEnv(scriptName string, env map[string]string) error {
	spec, ok := envScripts[scriptName]
	if !ok {
		return errx.Wrapf("validating env", fmt.Errorf("script %q is not an env script", scriptName), ErrorInvalidScript)
	}

	for _, required := range spec.requiredEnvs {
		if env[required] == "" {
			return errx.Wrapf("validating env", fmt.Errorf("%s is required by script %q", required, scriptName), ErrorInvalidEnv)
		}
	}

	known := map[string]bool{}
	for _, k := range spec.requiredEnvs {
		known[k] = true
	}
	for _, k := range spec.optionalEnvs {
		known[k] = true
	}

	for key := range env {
		if !envKeyPattern.MatchString(key) {
			return errx.Wrapf("validating env", fmt.Errorf("%q is not a valid env var name", key), ErrorInvalidEnv)
		}
		if !known[key] {
			return errx.Wrapf("validating env", fmt.Errorf("%q is not an allowed env var for script %q", key, scriptName), ErrorUnknownEnv)
		}
	}

	if validator, ok := scriptEnvValidators[scriptName]; ok {
		if err := validator(env); err != nil {
			return errx.Wrapf("validating env", err, ErrorInvalidEnv)
		}
	}

	return nil
}

// scripts whose env vars need a format check beyond "is it there"
var scriptEnvValidators = map[string]func(map[string]string) error{
	"ban-ip":     requiresValidIP,
	"unban-ip":   requiresValidIP,
	"ban-lookup": requiresValidIP,
}

func requiresValidIP(env map[string]string) error {
	if net.ParseIP(env["IPADDRESS"]) == nil {
		return fmt.Errorf("%q is not a valid IP address", env["IPADDRESS"])
	}
	return nil
}
