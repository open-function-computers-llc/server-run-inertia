package server

import (
	"net/url"
	"os"
	"strings"
)

const (
	cmsWordPress = "wordpress"
	cmsModx      = "modx"
)

// each server only hosts one CMS, and its hostname says which one, e.g.
// wp3.openfunction.website or modx1.openfunction.website
var cmsHostPrefixes = map[string]string{
	"wp":   cmsWordPress,
	"modx": cmsModx,
}

// detectCMS works out which CMS this server hosts. SERVER_CMS (wordpress|modx)
// overrides the hostname, which is handy for local dev on localhost.
func detectCMS(appURL string) string {
	if override := strings.ToLower(os.Getenv("SERVER_CMS")); override != "" {
		if override == cmsWordPress || override == cmsModx {
			return override
		}
		return ""
	}

	if !strings.Contains(appURL, "://") {
		appURL = "https://" + appURL
	}
	u, err := url.Parse(appURL)
	if err != nil {
		return ""
	}

	label, _, _ := strings.Cut(u.Hostname(), ".")
	prefix := strings.TrimRight(strings.ToLower(label), "0123456789")
	return cmsHostPrefixes[prefix]
}
