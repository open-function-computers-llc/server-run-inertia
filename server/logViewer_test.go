package server

import (
	"errors"
	"testing"

	"github.com/open-function-computers-llc/server-run-inertia/account"
)

func TestLogViewerArgs(t *testing.T) {
	t.Setenv("APACHE_LOG_DIR", "")
	accounts := []account.Account{{Name: "acme", PrimaryDomain: "acme.com"}}

	wp := &server{accounts: accounts, cms: cmsWordPress}
	got, err := wp.logViewerArgs(map[string]string{"LOG_TYPE": "access", "ACCOUNT_NAME": "acme"})
	if err != nil || got["LOG_FILE"] != "/var/log/httpd/acme_access.log" {
		t.Fatalf("wordpress: got %v %v", got, err)
	}

	modx := &server{accounts: accounts, cms: cmsModx}
	got, err = modx.logViewerArgs(map[string]string{"LOG_TYPE": "error", "ACCOUNT_NAME": "acme"})
	if err != nil || got["LOG_FILE"] != "/var/log/httpd/acme.com_error.log" {
		t.Fatalf("modx: got %v %v", got, err)
	}
}

func TestLogViewerArgsRejectsBadInput(t *testing.T) {
	s := &server{
		accounts: []account.Account{{Name: "acme", PrimaryDomain: "../../etc/passwd"}},
		cms:      cmsModx,
	}

	if _, err := s.logViewerArgs(map[string]string{"LOG_TYPE": "passwd", "ACCOUNT_NAME": "acme"}); !errors.Is(err, errorInvalidLogType) {
		t.Fatalf("expected invalid log type, got %v", err)
	}
	if _, err := s.logViewerArgs(map[string]string{"LOG_TYPE": "access", "ACCOUNT_NAME": "nobody"}); !errors.Is(err, ErrorAccountNotFound) {
		t.Fatalf("expected account not found, got %v", err)
	}
	if _, err := s.logViewerArgs(map[string]string{"LOG_TYPE": "access", "ACCOUNT_NAME": "acme"}); err == nil {
		t.Fatal("expected a path-like domain to be rejected")
	}
	// a LOG_FILE sent by the browser is ignored, not passed through
	s.cms = cmsWordPress
	got, _ := s.logViewerArgs(map[string]string{"LOG_TYPE": "access", "ACCOUNT_NAME": "acme", "LOG_FILE": "/etc/shadow"})
	if got["LOG_FILE"] != "/var/log/httpd/acme_access.log" {
		t.Fatalf("browser LOG_FILE leaked through: %v", got)
	}
}
