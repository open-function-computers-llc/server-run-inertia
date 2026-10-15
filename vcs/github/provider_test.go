package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/open-function-computers-llc/server-run-inertia/vcs"
)

func newTestProvider(t *testing.T, org string, handler http.HandlerFunc) gitHubProvider {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	t.Setenv("GITHUB_TOKEN", "test-token")
	t.Setenv("GITHUB_ORG", org)
	t.Setenv("GITHUB_API_URL", srv.URL)

	p, err := NewProvider()
	if err != nil {
		t.Fatal(err)
	}
	return p
}

func TestNewProviderRequiresToken(t *testing.T) {
	t.Setenv("GITHUB_TOKEN", "")
	if _, err := NewProvider(); !errors.Is(err, vcs.ErrorMissingConfig) {
		t.Fatalf("expected missing config, got %v", err)
	}
}

func TestListRepositoriesPaginatesOrgRepos(t *testing.T) {
	p := newTestProvider(t, "acme", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/orgs/acme/repos" {
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Errorf("unexpected auth header %q", got)
		}

		count := perPage
		if r.URL.Query().Get("page") == "2" {
			count = 1
		}
		repos := []map[string]string{}
		for i := 0; i < count; i++ {
			repos = append(repos, map[string]string{"full_name": fmt.Sprintf("acme/repo-%s-%d", r.URL.Query().Get("page"), i)})
		}
		json.NewEncoder(w).Encode(repos)
	})

	repos, err := p.ListRepositories()
	if err != nil {
		t.Fatal(err)
	}
	if len(repos) != perPage+1 {
		t.Fatalf("expected %d repos, got %d", perPage+1, len(repos))
	}
}

func TestListRepositoriesPersonal(t *testing.T) {
	p := newTestProvider(t, "", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/user/repos" || r.URL.Query().Get("affiliation") != "owner" {
			t.Errorf("unexpected request %s", r.URL)
		}
		fmt.Fprint(w, `[{"full_name":"me/site"}]`)
	})

	repos, err := p.ListRepositories()
	if err != nil || len(repos) != 1 || repos[0] != "me/site" {
		t.Fatalf("unexpected result %v %v", repos, err)
	}
}

func TestCreateRepositoryInOrg(t *testing.T) {
	p := newTestProvider(t, "acme", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/orgs/acme/repos" {
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
		}
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if body["name"] != "new-site" || body["private"] != true {
			t.Errorf("unexpected body %v", body)
		}
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{}`)
	})

	if err := p.CreateRepository("new-site"); err != nil {
		t.Fatal(err)
	}
}

func TestAddDeployKeyUsesFullName(t *testing.T) {
	p := newTestProvider(t, "acme", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/repos/other-owner/site/keys" {
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if body["title"] != "deploy" || body["read_only"] != true {
			t.Errorf("unexpected body %v", body)
		}
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{}`)
	})

	if err := p.AddDeployKey("other-owner/site", "deploy", "ssh-ed25519 AAAA", true); err != nil {
		t.Fatal(err)
	}
}

func TestAddDeployKeyLooksUpTokenOwner(t *testing.T) {
	p := newTestProvider(t, "", func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/user":
			fmt.Fprint(w, `{"login":"me"}`)
		case "/repos/me/site/keys":
			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, `{}`)
		default:
			t.Errorf("unexpected path %s", r.URL.Path)
		}
	})

	if err := p.AddDeployKey("site", "deploy", "ssh-ed25519 AAAA", false); err != nil {
		t.Fatal(err)
	}
}

func TestAPIErrorsAreReadable(t *testing.T) {
	p := newTestProvider(t, "acme", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprint(w, `{"message":"Validation Failed","errors":[{"message":"key is already in use"}]}`)
	})

	err := p.AddDeployKey("acme/site", "deploy", "ssh-ed25519 AAAA", true)
	if !errors.Is(err, vcs.ErrorProviderAPI) {
		t.Fatalf("expected provider API error, got %v", err)
	}
	if !strings.Contains(err.Error(), "key is already in use") || !strings.Contains(err.Error(), "422") {
		t.Fatalf("error should include GitHub's message, got %q", err)
	}
}
