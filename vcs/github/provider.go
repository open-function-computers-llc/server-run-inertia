package github

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/open-function-computers-llc/server-run-inertia/errx"
	"github.com/open-function-computers-llc/server-run-inertia/vcs"
)

const defaultAPIURL = "https://api.github.com"

type gitHubProvider struct {
	API_TOKEN string
	org       string
	apiURL    string
	client    *http.Client
}

// NewProvider reads GITHUB_TOKEN (required), GITHUB_ORG (optional, otherwise
// the token's own account is used) and GITHUB_API_URL (optional, for GitHub
// Enterprise).
func NewProvider() (gitHubProvider, error) {
	p := gitHubProvider{
		API_TOKEN: os.Getenv("GITHUB_TOKEN"),
		org:       os.Getenv("GITHUB_ORG"),
		apiURL:    strings.TrimRight(os.Getenv("GITHUB_API_URL"), "/"),
		client:    &http.Client{Timeout: 30 * time.Second},
	}

	if p.API_TOKEN == "" {
		return p, errx.Wrapf("configuring github", errors.New("GITHUB_TOKEN is not set"), vcs.ErrorMissingConfig)
	}
	if p.apiURL == "" {
		p.apiURL = defaultAPIURL
	}

	return p, nil
}

// do sends a request to the GitHub API and decodes a JSON response into out
// (if out isn't nil).
func (p gitHubProvider) do(method, path string, body, out any) error {
	var reqBody io.Reader
	if body != nil {
		encoded, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reqBody = bytes.NewReader(encoded)
	}

	req, err := http.NewRequest(method, p.apiURL+path, reqBody)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+p.API_TOKEN)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return errx.Wrapf("calling github", err, vcs.ErrorProviderAPI)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return errx.Wrapf("calling github", apiError(resp), vcs.ErrorProviderAPI)
	}

	if out == nil {
		return nil
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

// apiError turns a GitHub error response into a readable error, including the
// field-level details GitHub returns for validation failures.
func apiError(resp *http.Response) error {
	var payload struct {
		Message string `json:"message"`
		Errors  []struct {
			Message string `json:"message"`
			Field   string `json:"field"`
			Code    string `json:"code"`
		} `json:"errors"`
	}
	json.NewDecoder(resp.Body).Decode(&payload)

	msg := payload.Message
	if msg == "" {
		msg = resp.Status
	}
	for _, e := range payload.Errors {
		detail := e.Message
		if detail == "" {
			detail = strings.TrimSpace(e.Field + " " + e.Code)
		}
		if detail != "" {
			msg += ": " + detail
		}
	}

	return fmt.Errorf("%s (HTTP %d)", msg, resp.StatusCode)
}

// owner resolves who a repo belongs to: an explicit "owner/repo", then
// GITHUB_ORG, then the account the token belongs to.
func (p gitHubProvider) owner(repo string) (string, string, error) {
	if owner, name, found := strings.Cut(repo, "/"); found {
		return owner, name, nil
	}
	if p.org != "" {
		return p.org, repo, nil
	}

	var user struct {
		Login string `json:"login"`
	}
	if err := p.do(http.MethodGet, "/user", nil, &user); err != nil {
		return "", "", err
	}
	return user.Login, repo, nil
}
