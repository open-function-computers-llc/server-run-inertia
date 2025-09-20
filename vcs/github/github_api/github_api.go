package github_api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
)

func GitHubMakeRequest[T []map[string]any | map[string]any](url string, method string, body []byte) (T, error) {

	// Create the request
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(body))
	req.Header.Set("Authorization", "token "+os.Getenv("GITHUB_TOKEN"))
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Check if the request was successful
	if resp.StatusCode > 299 {
		// Parse the error response
		var errorResponse map[string]any
		err = json.Unmarshal(responseBody, &errorResponse)
		if err == nil && errorResponse["message"] != nil {
			return nil, errors.New("GitHub API Error: " + errorResponse["message"].(string))
		}
		// Generic error message
		return nil, errors.New("GitHub API Error: " + resp.Status)
	}

	// Parse the response body
	var responseBodyT T

	// Don't parse if the content-type is not application/json
	if !strings.Contains(resp.Header.Get("Content-Type"), "application/json") {
		// Return an empty instance of T
		return responseBodyT, nil
	}

	// JSON parse the response body
	err = json.Unmarshal(responseBody, &responseBodyT)
	if err != nil {
		return nil, err
	}

	return responseBodyT, nil
}
