package server

import (
	"net/http"
	"os/exec"
	"strings"

	scriptrunner "github.com/open-function-computers-llc/server-run-inertia/script-runner"
)

func (s *server) handleDashboard() http.HandlerFunc {
	// stuff to see if there are any accounts to import
	totalImportableAccounts := func() int {
		result, _ := scriptrunner.ReturnScriptOutput("list-importable-accounts", nil, nil)
		lines := strings.Split(string(result), "\n")
		totalImportableAccounts := 0
		for _, line := range lines {
			if len(strings.TrimSpace(line)) > 0 {
				totalImportableAccounts++
			}
		}

		return totalImportableAccounts
	}

	// stuff to show pie charts for the system disks
	type disc struct {
		Filesystem  string
		Size        string
		Available   string
		UsedPercent string
		MountPoint  string
	}
	allDiscs := []disc{}

	cmdOutput, _ := exec.Command("df", "-h").Output()
	lines := strings.Split(string(cmdOutput), "\n")
	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) < 6 {
			continue
		}

		if parts[5][0:1] != "/" {
			continue
		}

		d := disc{
			Filesystem:  parts[0],
			Size:        parts[1],
			Available:   parts[3],
			UsedPercent: strings.ReplaceAll(parts[4], "%", ""),
			MountPoint:  parts[5],
		}
		allDiscs = append(allDiscs, d)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		s.inertiaManager.Render(w, r, "Dashboard", map[string]any{
			"accounts":              s.accounts,
			"disks":                 allDiscs,
			"accountsCanBeImported": totalImportableAccounts() > 0,
		})
	}
}
