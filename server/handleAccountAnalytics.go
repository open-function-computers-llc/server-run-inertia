package server

import "net/http"

func (s *server) handleAccountAnalytics() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountName := r.PathValue("name")
		chartType := r.URL.Query().Get("type")

		bytes, status, _ := s.getAnalyticsJSON(accountName, chartType)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		w.Write(bytes)
	}
}
