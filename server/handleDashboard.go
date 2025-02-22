package server

import (
	"net/http"
	"os/exec"
	"strings"
)

func (s *server) handleDashboard() http.HandlerFunc {
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
			"totalAccounts": len(s.accounts),
			"disks":         allDiscs,
		})
	}
}
