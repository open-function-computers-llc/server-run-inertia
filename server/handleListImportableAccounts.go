package server

import (
	"net/http"
	"strings"

	scriptrunner "github.com/open-function-computers-llc/server-run-inertia/script-runner"
)

func (s *server) handleListImportableAccounts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, _ := scriptrunner.ReturnScriptOutput("list-importable-accounts", nil, nil)
		lines := strings.Split(string(result), "\n")
		importableAccounts := []string{}
		for _, line := range lines {
			if len(strings.TrimSpace(line)) > 0 {
				importableAccounts = append(importableAccounts, strings.TrimSpace(line))
			}
		}

		s.inertiaManager.Render(w, r, "Account/Importable", map[string]any{
			"accounts": importableAccounts,
		})
	}
}
