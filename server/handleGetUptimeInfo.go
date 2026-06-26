package server

import (
	"encoding/json"
	"net/http"
	"os"
	"time"
)

func (s *server) handleUptimeInfo() http.HandlerFunc {
	type incomingPayload struct {
		URL string `json:"url"`
	}

	type outputPayload struct {
		ErrorMessage string  `json:"error"`
		Day1         float64 `json:"day1"`
		Day7         float64 `json:"day7"`
		Day30        float64 `json:"day30"`
		Day60        float64 `json:"day60"`
		Day90        float64 `json:"day90"`
	}
	type uptimeResponsePayload struct {
		CertInfo struct {
			Valid   bool      `json:"valid"`
			Expires time.Time `json:"expires"`
			Error   string    `json:"error"`
			Names   []string  `json:"names"`
		} `json:"certInfo"`
		DomainInfo struct {
			Status          []string  `json:"status"`
			Registered      time.Time `json:"registered"`
			Expires         time.Time `json:"expires"`
			Registrant      string    `json:"registrant"`
			RegistrantEmail string    `json:"registrantEmail"`
			Registrar       string    `json:"registrar"`
			Error           string    `json:"error"`
		} `json:"domainInfo"`
		Outages []struct {
			ID       int       `json:"id"`
			SiteID   int       `json:"siteID"`
			Start    time.Time `json:"start"`
			End      time.Time `json:"end"`
			Duration int       `json:"duration"`
		} `json:"outages"`
		Up     bool `json:"up"`
		Uptime struct {
			Days1  float64 `json:"days1"`
			Days30 float64 `json:"days30"`
			Days60 float64 `json:"days60"`
			Days7  float64 `json:"days7"`
			Days90 float64 `json:"days90"`
		} `json:"uptime"`
		URL string `json:"url"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		uptimeUser := os.Getenv("UPTIME_USER")
		uptimePass := os.Getenv("UPTIME_PASS")
		uptimeDomain := os.Getenv("UPTIME_DOMAIN")

		var payload incomingPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			s.inertiaManager.Share("error", "Problem reading payload")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		output := outputPayload{}

		if payload.URL == "" {
			output.ErrorMessage = "Invalid URL: " + payload.URL
			sendJSONError(w, http.StatusBadRequest, output)
			return
		}

		resp, err := http.Get("https://" + uptimeUser + ":" + uptimePass + "@" + uptimeDomain + "/details?url=" + payload.URL)
		if err != nil {
			output.ErrorMessage = "Error communicating with the uptime provider"
			sendJSONError(w, http.StatusInternalServerError, output)
			return
		}
		defer resp.Body.Close()

		var uptimeResponse uptimeResponsePayload
		if err := json.NewDecoder(resp.Body).Decode(&uptimeResponse); err != nil {
			s.logger.Error("Uptime Error: " + err.Error())
			output.ErrorMessage = "Error understanding uptime provider output"
			sendJSONError(w, http.StatusInternalServerError, output)
			return
		}

		output.Day1 = uptimeResponse.Uptime.Days1
		output.Day7 = uptimeResponse.Uptime.Days7
		output.Day30 = uptimeResponse.Uptime.Days30
		output.Day60 = uptimeResponse.Uptime.Days60
		output.Day90 = uptimeResponse.Uptime.Days90

		sendJSON(w, output)
	}
}
