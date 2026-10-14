package server

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

const defaultApacheLogDir = "/var/log/httpd"

var errorInvalidLogType = errors.New("LOG_TYPE must be access or error")

// logViewerArgs turns the browser's request ({LOG_TYPE, ACCOUNT_NAME}) into the
// log file to tail. The path is always built here so the browser can't point
// tail at an arbitrary file.
//
// WordPress vhosts log to <account>_<type>.log, MODX vhosts to
// <domain>_<type>.log (see the ofco-deploy vhost templates).
func (s *server) logViewerArgs(args map[string]string) (map[string]string, error) {
	logType := args["LOG_TYPE"]
	if logType != "access" && logType != "error" {
		return nil, errorInvalidLogType
	}

	a, err := s.findAccount(args["ACCOUNT_NAME"])
	if err != nil {
		return nil, err
	}

	site := a.Name
	if s.cms == cmsModx {
		site = a.PrimaryDomain
	}
	if site == "" || strings.ContainsAny(site, `/\`) || strings.Contains(site, "..") {
		return nil, errors.New("can't work out the log file name for account " + a.Name)
	}

	logDir := os.Getenv("APACHE_LOG_DIR")
	if logDir == "" {
		logDir = defaultApacheLogDir
	}

	return map[string]string{
		"LOG_FILE": filepath.Join(logDir, site+"_"+logType+".log"),
	}, nil
}
