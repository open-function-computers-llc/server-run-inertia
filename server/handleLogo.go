package server

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/open-function-computers-llc/server-run-inertia/errx"
)

func (s *server) handleLogo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// default to the baked in OFC logo
		logoBytes := s.defaultLogo

		// check to see if we should load up an alternate logo
		customLogoBytes, err := useCustomLogo()
		if err == nil {
			logoBytes = customLogoBytes
		}

		output := map[string]string{
			"value": base64.StdEncoding.EncodeToString(logoBytes),
		}
		bytes, _ := json.Marshal(output)
		w.Write(bytes)
	}
}

func useCustomLogo() ([]byte, error) {
	filePath := "/etc/server-run-logo.png"

	info, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	if info.IsDir() {
		return nil, errx.Wrapf("custom logo path is a directory", fmt.Errorf("%s is a directory", filePath), ErrorLogoFileIsADirectory)
	}

	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return io.ReadAll(f)
}
