package server

import "errors"

var (
	ErrorAccountNotFound       = errors.New("account not found")
	ErrorSystemLoadOutputParse = errors.New("error parsing uptime command output")
	ErrorLogoFileIsADirectory  = errors.New("custom logo path is a directory")
)
