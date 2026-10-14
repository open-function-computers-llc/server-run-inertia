package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFlashErrorRoundTrip(t *testing.T) {
	s := &server{}

	// handler sets the flash and redirects
	setRec := httptest.NewRecorder()
	s.flashError(setRec, "Account not found: café/acme")
	cookies := setRec.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected 1 cookie, got %d", len(cookies))
	}

	// next request reads it once and clears it
	req := httptest.NewRequest(http.MethodGet, "/vcs", nil)
	req.AddCookie(cookies[0])
	readRec := httptest.NewRecorder()

	message, ok := s.consumeFlashError(readRec, req)
	if !ok || message != "Account not found: café/acme" {
		t.Fatalf("unexpected flash: %q %v", message, ok)
	}

	cleared := readRec.Result().Cookies()
	if len(cleared) != 1 || cleared[0].MaxAge >= 0 {
		t.Fatalf("expected flash cookie to be cleared, got %+v", cleared)
	}
}

func TestNoFlashError(t *testing.T) {
	s := &server{}

	req := httptest.NewRequest(http.MethodGet, "/vcs", nil)
	rec := httptest.NewRecorder()

	if _, ok := s.consumeFlashError(rec, req); ok {
		t.Fatal("expected no flash")
	}
	if len(rec.Result().Cookies()) != 0 {
		t.Fatal("should not touch cookies when there is no flash")
	}
}
