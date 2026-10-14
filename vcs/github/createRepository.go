package github

import (
	"net/http"
	"net/url"
)

func (p gitHubProvider) CreateRepository(name string) error {
	path := "/user/repos"
	if p.org != "" {
		path = "/orgs/" + url.PathEscape(p.org) + "/repos"
	}

	return p.do(http.MethodPost, path, map[string]any{
		"name":        name,
		"description": "This was created automatically from server-run",
		"private":     true,
		"auto_init":   true,
	}, nil)
}
