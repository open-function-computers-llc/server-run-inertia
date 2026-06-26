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

		bytes, status, err := s.getAnalyticsJSON(account, chartType)
		if err != nil {
			s.logger.Error(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		w.Write(bytes)
	}
}
