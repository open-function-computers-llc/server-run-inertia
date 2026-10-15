package github

import (
	"net/http"
	"net/url"
	"strconv"
)

const perPage = 100

func (p gitHubProvider) ListRepositories() ([]string, error) {
	path := "/user/repos?affiliation=owner&"
	if p.org != "" {
		path = "/orgs/" + url.PathEscape(p.org) + "/repos?"
	}

	output := []string{}
	for page := 1; ; page++ {
		var repos []struct {
			FullName string `json:"full_name"`
		}
		err := p.do(http.MethodGet, path+"per_page="+strconv.Itoa(perPage)+"&page="+strconv.Itoa(page), nil, &repos)
		if err != nil {
			return nil, err
		}

		for _, r := range repos {
			output = append(output, r.FullName)
		}
		if len(repos) < perPage {
			return output, nil
		}
	}
}
