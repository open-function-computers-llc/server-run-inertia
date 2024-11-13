package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
)

func (s *server) handleThirdPartyAnalyticsJSON() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		a := r.FormValue("account")
		t := r.FormValue("type")
		if a == "" || t == "" {
			sendJSONError(w, http.StatusBadRequest, map[string]string{
				"error": "the query params `account` and `type` are both required",
			})
			return
		}

		if !slices.Contains([]string{
			"bandwidth",
			"visitors",
		}, t) {
			sendJSONError(w, http.StatusBadRequest, map[string]string{
				"error": "invalid `type`, pass either 'bandwidth' or 'visitors'",
			})
			return
		}

		fmt.Println("got past the checks", a, t)

		chartData, err := processAnalyticDataForAccount(a, t)
		if err != nil {
			sendJSONError(w, http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
			return
		}

		data := map[string]thirdPartyAnalyticsJSON{
			t: {
				Values:   chartData.outputValues,
				Duration: len(chartData.outputValues),
				Start:    chartData.outputLabels[0],
				End:      chartData.outputLabels[len(chartData.outputLabels)-1],
				Average:  getAverage(chartData.outputValues),
				Total:    getTotal(chartData.outputValues),
				Min:      0,
				Max:      getMax(chartData.outputValues),
			},
		}
		fmt.Println(data)

		bytes, _ := json.Marshal(data)
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
