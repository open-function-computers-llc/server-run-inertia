package server

import "testing"

func TestDetectCMS(t *testing.T) {
	cases := map[string]string{
		"https://wp3.openfunction.website":        cmsWordPress,
		"https://modx1.openfunction.website:8420": cmsModx,
		"wp.openfunction.website":                 cmsWordPress,
		"https://WP12.openfunction.website/":      cmsWordPress,
		"http://localhost:8420":                   "",
		"https://wpengine.example.com":            "",
		"https://openfunction.website":            "",
	}

	for in, want := range cases {
		if got := detectCMS(in); got != want {
			t.Errorf("detectCMS(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestDetectCMSOverride(t *testing.T) {
	t.Setenv("SERVER_CMS", "MODX")
	if got := detectCMS("http://localhost:8420"); got != cmsModx {
		t.Fatalf("expected override to win, got %q", got)
	}

	t.Setenv("SERVER_CMS", "drupal")
	if got := detectCMS("https://wp3.openfunction.website"); got != "" {
		t.Fatalf("expected unknown override to disable, got %q", got)
	}
}
