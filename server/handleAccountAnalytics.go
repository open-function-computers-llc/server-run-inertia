package server

import "net/http"

func (s *server) handleAccountAnalytics() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		account := r.PathValue("account")
		if account == "" {
			account = r.URL.Query().Get("account")
		}

		chartType := r.URL.Query().Get("type")
		if chartType == "" {
			chartType = "bandwidth"
		}

		bytes, status, _ := s.getAnalyticsJSON(account, chartType)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		w.Write(bytes)
	}
}
