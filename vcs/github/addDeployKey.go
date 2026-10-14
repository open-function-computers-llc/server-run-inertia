package github

import (
	"net/http"
	"net/url"
)

func (p gitHubProvider) AddDeployKey(repo, name, key string, readOnly bool) error {
	owner, repo, err := p.owner(repo)
	if err != nil {
		return err
	}

	return p.do(http.MethodPost, "/repos/"+url.PathEscape(owner)+"/"+url.PathEscape(repo)+"/keys", map[string]any{
		"title":     name,
		"key":       key,
		"read_only": readOnly,
	}, nil)
}
