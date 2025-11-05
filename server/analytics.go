package server

import (
	"encoding/json"
	"net/http"
	"os"
	"slices"
	"strings"
)

type thirdPartyAnalyticsJSON struct {
	Values   []float64 `json:"values"`
	Min      float64   `json:"min"`
	Max      float64   `json:"max"`
	Start    string    `json:"start"`
	End      string    `json:"end"`
	Average  float64   `json:"average"`
	Total    float64   `json:"total"`
	Duration int       `json:"duration"`
}

type anaylticJSON struct {
	General struct {
		StartDate         string   `json:"start_date"`
		EndDate           string   `json:"end_date"`
		DateTime          string   `json:"date_time"`
		TotalRequests     int      `json:"total_requests"`
		ValidRequests     int      `json:"valid_requests"`
		FailedRequests    int      `json:"failed_requests"`
		GenerationTime    int      `json:"generation_time"`
		UniqueVisitors    int      `json:"unique_visitors"`
		UniqueFiles       int      `json:"unique_files"`
		ExcludedHits      int      `json:"excluded_hits"`
		UniqueReferrers   int      `json:"unique_referrers"`
		UniqueNotFound    int      `json:"unique_not_found"`
		UniqueStaticFiles int      `json:"unique_static_files"`
		LogSize           int      `json:"log_size"`
		Bandwidth         int      `json:"bandwidth"`
		LogPath           []string `json:"log_path"`
	} `json:"general"`
	Visitors struct {
		Metadata struct {
			Bytes struct {
				Count int `json:"count"`
				Max   int `json:"max"`
				Min   int `json:"min"`
			} `json:"bytes"`
			Visitors struct {
				Count int `json:"count"`
				Max   int `json:"max"`
				Min   int `json:"min"`
			} `json:"visitors"`
			Hits struct {
				Count int `json:"count"`
				Max   int `json:"max"`
				Min   int `json:"min"`
			} `json:"hits"`
			Data struct {
				Unique int `json:"unique"`
			} `json:"data"`
		} `json:"metadata"`
		Data []struct {
			Hits struct {
				Count   int    `json:"count"`
				Percent string `json:"percent"`
			} `json:"hits"`
			Visitors struct {
				Count   int    `json:"count"`
				Percent string `json:"percent"`
			} `json:"visitors"`
			Bytes struct {
				Count   int    `json:"count"`
				Percent string `json:"percent"`
			} `json:"bytes"`
			Data string `json:"data"`
		} `json:"data"`
	} `json:"visitors"`
	Requests struct {
		Metadata struct {
			Bytes struct {
				Count int `json:"count"`
				Max   int `json:"max"`
				Min   int `json:"min"`
			} `json:"bytes"`
			Visitors struct {
				Count int `json:"count"`
				Max   int `json:"max"`
				Min   int `json:"min"`
			} `json:"visitors"`
			Hits struct {
				Count int `json:"count"`
				Max   int `json:"max"`
				Min   int `json:"min"`
			} `json:"hits"`
			Data struct {
				Unique int `json:"unique"`
			} `json:"data"`
		} `json:"metadata"`
		Data []struct {
			Hits struct {
				Count   int    `json:"count"`
				Percent string `json:"percent"`
			} `json:"hits"`
			Visitors struct {
				Count   int    `json:"count"`
				Percent string `json:"percent"`
			} `json:"visitors"`
			Bytes struct {
				Count   int    `json:"count"`
				Percent string `json:"percent"`
			} `json:"bytes"`
			Method   string `json:"method"`
			Protocol string `json:"protocol"`
			Data     string `json:"data"`
		} `json:"data"`
	} `json:"requests"`
	StaticRequests struct {
		Metadata struct {
			Bytes struct {
				Count int `json:"count"`
				Max   int `json:"max"`
				Min   int `json:"min"`
			} `json:"bytes"`
			Visitors struct {
				Count int `json:"count"`
				Max   int `json:"max"`
				Min   int `json:"min"`
			} `json:"visitors"`
			Hits struct {
				Count int `json:"count"`
				Max   int `json:"max"`
				Min   int `json:"min"`
			} `json:"hits"`
			Data struct {
				Unique int `json:"unique"`
			} `json:"data"`
		} `json:"metadata"`
		Data []struct {
			Hits struct {
				Count   int    `json:"count"`
				Percent string `json:"percent"`
			} `json:"hits"`
			Visitors struct {
				Count   int    `json:"count"`
				Percent string `json:"percent"`
			} `json:"visitors"`
			Bytes struct {
				Count   int    `json:"count"`
				Percent string `json:"percent"`
			} `json:"bytes"`
			Method   string `json:"method"`
			Protocol string `json:"protocol"`
			Data     string `json:"data"`
		} `json:"data"`
	} `json:"static_requests"`
	NotFound struct {
		Metadata struct {
			Bytes struct {
				Count int `json:"count"`
				Max   int `json:"max"`
				Min   int `json:"min"`
			} `json:"bytes"`
			Visitors struct {
				Count int `json:"count"`
				Max   int `json:"max"`
				Min   int `json:"min"`
			} `json:"visitors"`
			Hits struct {
				Count int `json:"count"`
				Max   int `json:"max"`
				Min   int `json:"min"`
			} `json:"hits"`
			Data struct {
				Unique int `json:"unique"`
			} `json:"data"`
		} `json:"metadata"`
		Data []struct {
			Hits struct {
				Count   int    `json:"count"`
				Percent string `json:"percent"`
			} `json:"hits"`
			Visitors struct {
				Count   int    `json:"count"`
				Percent string `json:"percent"`
			} `json:"visitors"`
			Bytes struct {
				Count   int    `json:"count"`
				Percent string `json:"percent"`
			} `json:"bytes"`
			Method   string `json:"method"`
			Protocol string `json:"protocol"`
			Data     string `json:"data"`
		} `json:"data"`
	} `json:"not_found"`
	Hosts struct {
		Metadata struct {
			Bytes struct {
				Count int `json:"count"`
				Max   int `json:"max"`
				Min   int `json:"min"`
			} `json:"bytes"`
			Visitors struct {
				Count int `json:"count"`
				Max   int `json:"max"`
				Min   int `json:"min"`
			} `json:"visitors"`
			Hits struct {
				Count int `json:"count"`
				Max   int `json:"max"`
				Min   int `json:"min"`
			} `json:"hits"`
			Data struct {
				Unique int `json:"unique"`
			} `json:"data"`
		} `json:"metadata"`
		Data []struct {
			Hits struct {
				Count   int    `json:"count"`
				Percent string `json:"percent"`
			} `json:"hits"`
			Visitors struct {
				Count   int    `json:"count"`
				Percent string `json:"percent"`
			} `json:"visitors"`
			Bytes struct {
				Count   int    `json:"count"`
				Percent string `json:"percent"`
			} `json:"bytes"`
			Data    string   `json:"data"`
			Country string   `json:"country"`
			Items   []string `json:"items"`
		} `json:"data"`
	} `json:"hosts"`
	Os struct {
		Metadata struct {
			Bytes struct {
				Count int `json:"count"`
				Max   int `json:"max"`
				Min   int `json:"min"`
			} `json:"bytes"`
			Visitors struct {
				Count int `json:"count"`
				Max   int `json:"max"`
				Min   int `json:"min"`
			} `json:"visitors"`
			Hits struct {
				Count int `json:"count"`
				Max   int `json:"max"`
				Min   int `json:"min"`
			} `json:"hits"`
			Data struct {
				Unique int `json:"unique"`
			} `json:"data"`
		} `json:"metadata"`
		Data []struct {
			Hits struct {
				Count   int    `json:"count"`
				Percent string `json:"percent"`
			} `json:"hits"`
			Visitors struct {
				Count   int    `json:"count"`
				Percent string `json:"percent"`
			} `json:"visitors"`
			Bytes struct {
				Count   int    `json:"count"`
				Percent string `json:"percent"`
			} `json:"bytes"`
			Data  string `json:"data"`
			Items []struct {
				Hits struct {
					Count   int    `json:"count"`
					Percent string `json:"percent"`
				} `json:"hits"`
				Visitors struct {
					Count   int    `json:"count"`
					Percent string `json:"percent"`
				} `json:"visitors"`
				Bytes struct {
					Count   int    `json:"count"`
					Percent string `json:"percent"`
				} `json:"bytes"`
				Data string `json:"data"`
			} `json:"items"`
		} `json:"data"`
	} `json:"os"`
	Browsers struct {
		Metadata struct {
			Bytes struct {
				Count int `json:"count"`
				Max   int `json:"max"`
				Min   int `json:"min"`
			} `json:"bytes"`
			Visitors struct {
				Count int `json:"count"`
				Max   int `json:"max"`
				Min   int `json:"min"`
			} `json:"visitors"`
			Hits struct {
				Count int `json:"count"`
				Max   int `json:"max"`
				Min   int `json:"min"`
			} `json:"hits"`
			Data struct {
				Unique int `json:"unique"`
			} `json:"data"`
		} `json:"metadata"`
		Data []struct {
			Hits struct {
				Count   int    `json:"count"`
				Percent string `json:"percent"`
			} `json:"hits"`
			Visitors struct {
				Count   int    `json:"count"`
				Percent string `json:"percent"`
			} `json:"visitors"`
			Bytes struct {
				Count   int    `json:"count"`
				Percent string `json:"percent"`
			} `json:"bytes"`
			Data  string `json:"data"`
			Items []struct {
				Hits struct {
					Count   int    `json:"count"`
					Percent string `json:"percent"`
				} `json:"hits"`
				Visitors struct {
					Count   int    `json:"count"`
					Percent string `json:"percent"`
				} `json:"visitors"`
				Bytes struct {
					Count   int    `json:"count"`
					Percent string `json:"percent"`
				} `json:"bytes"`
				Data string `json:"data"`
			} `json:"items"`
		} `json:"data"`
	} `json:"browsers"`
	VisitTime struct {
		Metadata struct {
			Bytes struct {
				Count int `json:"count"`
				Max   int `json:"max"`
				Min   int `json:"min"`
			} `json:"bytes"`
			Visitors struct {
				Count int `json:"count"`
				Max   int `json:"max"`
				Min   int `json:"min"`
			} `json:"visitors"`
			Hits struct {
				Count int `json:"count"`
				Max   int `json:"max"`
				Min   int `json:"min"`
			} `json:"hits"`
			Data struct {
				Unique int `json:"unique"`
			} `json:"data"`
		} `json:"metadata"`
		Data []struct {
			Hits struct {
				Count   int    `json:"count"`
				Percent string `json:"percent"`
			} `json:"hits"`
			Visitors struct {
				Count   int    `json:"count"`
				Percent string `json:"percent"`
			} `json:"visitors"`
			Bytes struct {
				Count   int    `json:"count"`
				Percent string `json:"percent"`
			} `json:"bytes"`
			Data string `json:"data"`
		} `json:"data"`
	} `json:"visit_time"`
	ReferringSites struct {
		Metadata struct {
			Bytes struct {
				Count int `json:"count"`
				Max   int `json:"max"`
				Min   int `json:"min"`
			} `json:"bytes"`
			Visitors struct {
				Count int `json:"count"`
				Max   int `json:"max"`
				Min   int `json:"min"`
			} `json:"visitors"`
			Hits struct {
				Count int `json:"count"`
				Max   int `json:"max"`
				Min   int `json:"min"`
			} `json:"hits"`
			Data struct {
				Unique int `json:"unique"`
			} `json:"data"`
		} `json:"metadata"`
		Data []struct {
			Hits struct {
				Count   int    `json:"count"`
				Percent string `json:"percent"`
			} `json:"hits"`
			Visitors struct {
				Count   int    `json:"count"`
				Percent string `json:"percent"`
			} `json:"visitors"`
			Bytes struct {
				Count   int    `json:"count"`
				Percent string `json:"percent"`
			} `json:"bytes"`
			Data string `json:"data"`
		} `json:"data"`
	} `json:"referring_sites"`
	StatusCodes struct {
		Metadata struct {
			Bytes struct {
				Count int `json:"count"`
				Max   int `json:"max"`
				Min   int `json:"min"`
			} `json:"bytes"`
			Visitors struct {
				Count int `json:"count"`
				Max   int `json:"max"`
				Min   int `json:"min"`
			} `json:"visitors"`
			Hits struct {
				Count int `json:"count"`
				Max   int `json:"max"`
				Min   int `json:"min"`
			} `json:"hits"`
			Data struct {
				Unique int `json:"unique"`
			} `json:"data"`
		} `json:"metadata"`
		Data []struct {
			Hits struct {
				Count   int    `json:"count"`
				Percent string `json:"percent"`
			} `json:"hits"`
			Visitors struct {
				Count   int    `json:"count"`
				Percent string `json:"percent"`
			} `json:"visitors"`
			Bytes struct {
				Count   int    `json:"count"`
				Percent string `json:"percent"`
			} `json:"bytes"`
			Data  string `json:"data"`
			Items []struct {
				Hits struct {
					Count   int    `json:"count"`
					Percent string `json:"percent"`
				} `json:"hits"`
				Visitors struct {
					Count   int    `json:"count"`
					Percent string `json:"percent"`
				} `json:"visitors"`
				Bytes struct {
					Count   int    `json:"count"`
					Percent string `json:"percent"`
				} `json:"bytes"`
				Data string `json:"data"`
			} `json:"items"`
		} `json:"data"`
	} `json:"status_codes"`
	Geolocation struct {
		Metadata struct {
			Bytes struct {
				Count int `json:"count"`
				Max   int `json:"max"`
				Min   int `json:"min"`
			} `json:"bytes"`
			Visitors struct {
				Count int `json:"count"`
				Max   int `json:"max"`
				Min   int `json:"min"`
			} `json:"visitors"`
			Hits struct {
				Count int `json:"count"`
				Max   int `json:"max"`
				Min   int `json:"min"`
			} `json:"hits"`
			Data struct {
				Unique int `json:"unique"`
			} `json:"data"`
		} `json:"metadata"`
		Data []struct {
			Hits struct {
				Count   int    `json:"count"`
				Percent string `json:"percent"`
			} `json:"hits"`
			Visitors struct {
				Count   int    `json:"count"`
				Percent string `json:"percent"`
			} `json:"visitors"`
			Bytes struct {
				Count   int    `json:"count"`
				Percent string `json:"percent"`
			} `json:"bytes"`
			Data  string `json:"data"`
			Items []struct {
				Hits struct {
					Count   int    `json:"count"`
					Percent string `json:"percent"`
				} `json:"hits"`
				Visitors struct {
					Count   int    `json:"count"`
					Percent string `json:"percent"`
				} `json:"visitors"`
				Bytes struct {
					Count   int    `json:"count"`
					Percent string `json:"percent"`
				} `json:"bytes"`
				Data string `json:"data"`
			} `json:"items"`
		} `json:"data"`
	} `json:"geolocation"`
}

type processedChartData struct {
	outputValues     []float64
	outputLabels     []string
	title            string
	chartType        string
	yAxisLabelFormat string
	seriesName       string
	tooltipSuffix    string
	tooltipPrefix    string
	tooltipDecimals  int
}

func processAnalyticDataForAccount(account string, dataToFetch string) (processedChartData, error) {
	chartData := processedChartData{
		outputValues:     []float64{},
		outputLabels:     []string{},
		title:            "",
		chartType:        "column",
		yAxisLabelFormat: "{value:,.0f}",
		seriesName:       "Data",
		tooltipSuffix:    "",
		tooltipPrefix:    "",
		tooltipDecimals:  0,
	}

	// TODO: look in reports folder for all json files that start with the account name
	allReports, err := os.ReadDir(os.Getenv("REPORTS_ROOT"))
	if err != nil {
		return chartData, err
	}

	for _, report := range allReports {
		_, err := report.Info()
		if err != nil {
			return chartData, err
		}

		if strings.HasPrefix(report.Name(), account) {
			reportBytes, _ := os.ReadFile(os.Getenv("REPORTS_ROOT") + report.Name())
			var reportJSON anaylticJSON
			err := json.Unmarshal(reportBytes, &reportJSON)
			if err != nil {
				continue
			}

			if dataToFetch == "total-requests" {
				chartData.outputValues = append(chartData.outputValues, float64(reportJSON.General.TotalRequests))
				chartData.title = "Total Requests"
				chartData.seriesName = "Requests"
			}

			if dataToFetch == "unique-visitors" || dataToFetch == "visitors" {
				chartData.outputValues = append(chartData.outputValues, float64(reportJSON.General.UniqueVisitors))
				chartData.title = "Unique Visitors"
				chartData.seriesName = "Visitors"
			}

			if dataToFetch == "bandwidth" {
				chartData.outputValues = append(chartData.outputValues, float64(reportJSON.General.Bandwidth)/1000000)
				chartData.title = "Bandwidth"
				chartData.chartType = "line"
				chartData.yAxisLabelFormat = "{value:,.0f} MB"
				chartData.seriesName = "Bandwidth"
				chartData.tooltipSuffix = " MB"
				chartData.tooltipDecimals = 2
			}
			chartData.outputLabels = append(chartData.outputLabels, reportJSON.General.EndDate)
		}
	}

	return chartData, nil
}

func (s *server) getAnalyticsJSON(accountName, chartType string) ([]byte, int, error) {
	if accountName == "" || chartType == "" {
		errResp, _ := json.Marshal(map[string]string{
			"error": "the query params `account` and `type` are both required",
		})
		return errResp, http.StatusBadRequest, nil
	}

	if !slices.Contains([]string{"bandwidth", "unique-visitors, total-requests"}, chartType) {
		errResp, _ := json.Marshal(map[string]string{
			"error": "invalid `type`, pass either 'bandwidth' or 'visitors'",
		})
		return errResp, http.StatusBadRequest, nil
	}

	chartData, err := processAnalyticDataForAccount(accountName, chartType)
	if err != nil {
		errResp, _ := json.Marshal(map[string]string{
			"error": err.Error(),
		})
		return errResp, http.StatusBadRequest, nil
	}

	if len(chartData.outputLabels) == 0 || len(chartData.outputValues) == 0 {
		errResp, _ := json.Marshal(map[string]string{
			"error": "no analytics data found for this account",
		})
		return errResp, http.StatusNotFound, nil
	}

	data := map[string]thirdPartyAnalyticsJSON{
		chartType: {
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

	bytes, _ := json.Marshal(data)
	return bytes, http.StatusOK, nil
}
