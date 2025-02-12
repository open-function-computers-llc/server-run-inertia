package server

import (
	"errors"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

func (s *server) handleAPISystemLoad() http.HandlerFunc {
	type stats struct {
		OneMinute      string `json:"oneMinute"`
		FiveMinutes    string `json:"fiveMinutes"`
		FifteenMinutes string `json:"fifteenMinutes"`
	}

	parseMac := func(stdout string) (stats, error) {
		output := stats{}
		parts := strings.Split(stdout, "load averages: ")
		if len(parts) != 2 {
			output.OneMinute = "error"
			output.FiveMinutes = "error"
			output.FifteenMinutes = "error"
			return output, nil
		}

		loadParts := strings.Split(parts[1], " ")
		output.OneMinute = loadParts[0]
		output.FiveMinutes = loadParts[1]
		output.FifteenMinutes = loadParts[2]

		return output, nil
	}

	parse := func(stdout []byte) (stats, error) {
		output := stats{}
		stdOutString := string(stdout)
		parts := strings.Split(stdOutString, "load average: ")
		if len(parts) != 2 {
			return output, errors.New(SystemLoadOutputParse)
		}

		loadParts := strings.Split(parts[1], ", ")
		output.OneMinute = loadParts[0]
		output.FiveMinutes = loadParts[1]
		output.FifteenMinutes = strings.TrimSpace(loadParts[2])

		return output, nil
	}

	return func(w http.ResponseWriter, r *http.Request) {
		c, err := s.upgrader.Upgrade(w, r, nil)

		if err != nil {
			sendJSONError(w, http.StatusInternalServerError, map[string]string{
				"error": "Problem setting up the web socket connection",
			})
		}
		defer c.Close()

		// start streaming the system load metrics
		for {
			stdout, err := exec.Command("uptime").Output()
			if err != nil {
				c.WriteJSON(map[string]string{"error": err.Error()})
			}

			stats, err := parse(stdout)
			if err != nil {
				// try to parse it from mac output?
				stats, err = parseMac(string(stdout))
				if err != nil {
					c.WriteJSON(map[string]string{"error": err.Error()})
				}
			}

			c.WriteJSON(stats)
			time.Sleep(5 * time.Second)
		}
	}
}
