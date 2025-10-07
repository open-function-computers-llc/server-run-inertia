package server

import (
	"net/http"
)

func (s *server) handleThirdPartyAnalyticsJSON() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		accountName := r.FormValue("account")
		chartType := r.FormValue("type")

		bytes, status, _ := s.getAnalyticsJSON(accountName, chartType)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		w.Write(bytes)
	}

}

func getAverage(data []float64) float64 {
	total := 0.0
	for i := 0; i < len(data); i++ {
		total += data[i]
	}
	return total / float64(len(data))
}

func getTotal(data []float64) float64 {
	total := 0.0
	for i := 0; i < len(data); i++ {
		total += data[i]
	}
	return total
}

func getMax(data []float64) float64 {
	highest := 0.0
	for i := 0; i < len(data); i++ {
		if data[i] > highest {
			highest = data[i]
		}
	}
	return highest
}
