package server

import (
	"fmt"
	"net/http"
	"os/exec"
	"strconv"
	"strings"

	scriptrunner "github.com/open-function-computers-llc/server-run-inertia/script-runner"
)

func (s *server) handleDashboard() http.HandlerFunc {
	// stuff to calculate the fail to ban jails
	type jail struct {
		Name            string
		TotalBlockedIPs int
	}

	dashboardFailToBanStats := func() []jail {
		output := []jail{}
		result, _ := scriptrunner.ReturnScriptOutput("f2banstatus", nil, nil)

		lines := strings.Split(string(result), "\n")
		validNames := []string{}
		validTotals := []int{}
		for _, line := range lines {
			if len(line) < 17 {
				continue
			}

			if line[:16] == "::::::::::::::::" {
				name := strings.ReplaceAll(line, ":", "")
				validNames = append(validNames, strings.TrimSpace(name))
			}
			if line[:17] == "|- Total banned: " {
				totalAsString := strings.ReplaceAll(line, "|- Total banned: ", "")
				total, err := strconv.Atoi(totalAsString)
				if err != nil {
					total = -1
				}
				validTotals = append(validTotals, total)
			}
		}

		if len(validNames) != len(validTotals) {
			errorJail := jail{
				Name: "There was an error getting the fail2ban info on this server",
			}
			output = append(output, errorJail)
			return output
		}

		for i := range validNames {
			jail := jail{
				Name:            validNames[i],
				TotalBlockedIPs: validTotals[i],
			}
			output = append(output, jail)
		}
		fmt.Println(output)

		return output
	}

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
			"f2bJails":              dashboardFailToBanStats(),
			"accountsCanBeImported": totalImportableAccounts() > 0,
		})
	}
}
