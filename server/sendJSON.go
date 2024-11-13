package server

import (
	"encoding/json"
	"net/http"
)

func sendJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Add("content-type", "application/json")
	bytes, err := json.Marshal(data)
	if err != nil {
		sendJSONError(w, http.StatusInternalServerError, map[string]string{
			"error": "Internal system error. Sorry bout that!",
		})
	}

	w.Write(bytes)
}

func sendJSONError(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(statusCode)
	bytes, _ := json.Marshal(data)

	w.Write(bytes)
}
