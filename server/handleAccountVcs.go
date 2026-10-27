package server

import (
	"encoding/json"
	"net/http"

	"github.com/open-function-computers-llc/server-run-inertia/account"
	"github.com/open-function-computers-llc/server-run-inertia/vcs/registry"
)

// VCS is per client: each handler here resolves the provider from the
// account's own setting (falling back to the server default).

func (s *server) handleAccountVcsRepos() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a, err := account.Find(r.PathValue("account"))
		if err != nil {
			sendJSONError(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}

		p, name, err := registry.ForAccount(a.VCSProvider)
		if err != nil {
			sendJSONError(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		repos, err := p.ListRepositories()
		if err != nil {
			sendJSONError(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
			return
		}

		sendJSON(w, map[string]any{
			"provider": name,
			"repos":    repos,
		})
	}
}

func (s *server) handleAccountVcsCreateRepo() http.HandlerFunc {
	type incomingPayload struct {
		Name string `json:"name"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		a, err := account.Find(r.PathValue("account"))
		if err != nil {
			s.flashError(w, err.Error())
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			return
		}

		var payload incomingPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload.Name == "" {
			s.flashError(w, "Repository name is required")
			http.Redirect(w, r, "/account/"+a.Name, http.StatusSeeOther)
			return
		}

		p, _, err := registry.ForAccount(a.VCSProvider)
		if err != nil {
			s.flashError(w, err.Error())
			http.Redirect(w, r, "/account/"+a.Name, http.StatusSeeOther)
			return
		}

		if err := p.CreateRepository(payload.Name); err != nil {
			s.flashError(w, err.Error())
		}

		http.Redirect(w, r, "/account/"+a.Name, http.StatusSeeOther)
	}
}

func (s *server) handleAccountVcsAddDeployKey() http.HandlerFunc {
	type incomingPayload struct {
		Repo     string `json:"repo"`
		Name     string `json:"name"`
		Key      string `json:"key"`
		ReadOnly bool   `json:"readOnly"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		a, err := account.Find(r.PathValue("account"))
		if err != nil {
			s.flashError(w, err.Error())
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			return
		}

		var payload incomingPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			s.flashError(w, "Problem reading payload")
			http.Redirect(w, r, "/account/"+a.Name, http.StatusSeeOther)
			return
		}

		if payload.Repo == "" || payload.Name == "" || payload.Key == "" {
			s.flashError(w, "Repo, name, and key are all required")
			http.Redirect(w, r, "/account/"+a.Name, http.StatusSeeOther)
			return
		}

		p, _, err := registry.ForAccount(a.VCSProvider)
		if err != nil {
			s.flashError(w, err.Error())
			http.Redirect(w, r, "/account/"+a.Name, http.StatusSeeOther)
			return
		}

		if err := p.AddDeployKey(payload.Repo, payload.Name, payload.Key, payload.ReadOnly); err != nil {
			s.flashError(w, err.Error())
		}

		http.Redirect(w, r, "/account/"+a.Name, http.StatusSeeOther)
	}
}

func (s *server) handleAccountVcsSetProvider() http.HandlerFunc {
	type incomingPayload struct {
		Provider string `json:"provider"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		a, err := account.Find(r.PathValue("account"))
		if err != nil {
			s.flashError(w, err.Error())
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			return
		}

		var payload incomingPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			s.flashError(w, "Problem reading payload")
			http.Redirect(w, r, "/account/"+a.Name, http.StatusSeeOther)
			return
		}

		if err := a.SetVCSProvider(payload.Provider); err != nil {
			s.flashError(w, err.Error())
		}

		http.Redirect(w, r, "/account/"+a.Name, http.StatusSeeOther)
	}
}
