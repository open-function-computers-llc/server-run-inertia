package scriptrunner

import (
	"errors"
	"testing"
)

func TestValidateEnvRequiredVars(t *testing.T) {
	err := validateEnv("ban-ip", map[string]string{})
	if !errors.Is(err, ErrorInvalidEnv) {
		t.Fatalf("expected missing env error, got %v", err)
	}

	err = validateEnv("ban-ip", map[string]string{"IPADDRESS": ""})
	if !errors.Is(err, ErrorInvalidEnv) {
		t.Fatalf("expected missing env error for empty value, got %v", err)
	}
}

func TestValidateEnvRejectsUnregisteredKeys(t *testing.T) {
	err := validateEnv("add-account", map[string]string{"ACCOUNT_NAME": "blog", "SNEAKY": "yes"})
	if !errors.Is(err, ErrorUnknownEnv) {
		t.Fatalf("expected unknown env error, got %v", err)
	}

	err = validateEnv("add-account", map[string]string{"ACCOUNT_NAME": "blog", "lowercase": "yes"})
	if !errors.Is(err, ErrorInvalidEnv) {
		t.Fatalf("expected invalid env name error, got %v", err)
	}
}

func TestValidateEnvAllowsRegisteredOptionals(t *testing.T) {
	err := validateEnv("add-account", map[string]string{
		"ACCOUNT_NAME":   "blog",
		"INSTALL_WP":     "true",
		"INSTALL_RAD":    "false",
		"ADMIN_PASSWORD": "p@ss&w0rd#%+",
		"SITE_TITLE":     "My Site",
		"ADMIN_USERNAME": "admin",
		"ADMIN_EMAIL":    "admin@example.com",
		"DOMAIN":         "example.com",
	})
	if err != nil {
		t.Fatalf("expected valid env, got %v", err)
	}
}

func TestValidateEnvIPFormat(t *testing.T) {
	for _, good := range []string{"203.0.113.7", "2001:db8::1"} {
		if err := validateEnv("ban-ip", map[string]string{"IPADDRESS": good}); err != nil {
			t.Fatalf("%q should be valid, got %v", good, err)
		}
	}

	for _, bad := range []string{"203.0.113.999", "not-an-ip", "203.0.113", "555555"} {
		err := validateEnv("ban-lookup", map[string]string{"IPADDRESS": bad})
		if !errors.Is(err, ErrorInvalidEnv) {
			t.Fatalf("%q should be invalid, got %v", bad, err)
		}
	}
}

func TestValidateEnvUnknownScript(t *testing.T) {
	err := validateEnv("not-a-script", map[string]string{})
	if !errors.Is(err, ErrorInvalidScript) {
		t.Fatalf("expected invalid script error, got %v", err)
	}
}

func TestClassifyLine(t *testing.T) {
	if m := classifyLine("##section## Create WordPress"); m.Type != MsgSection || m.Message != "Create WordPress" {
		t.Fatalf("section line misclassified: %+v", m)
	}
	if m := classifyLine("##error## boom"); m.Type != MsgError || m.Message != "boom" {
		t.Fatalf("error line misclassified: %+v", m)
	}
	if m := classifyLine("just output"); m.Type != MsgStdout || m.Message != "just output" {
		t.Fatalf("stdout line misclassified: %+v", m)
	}
}
