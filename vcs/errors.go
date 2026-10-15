package vcs

import "errors"

var (
	ErrorNoProviderConfigured = errors.New("no VCS provider configured")
	ErrorNotImplemented       = errors.New("not yet implemented for this VCS provider")
	ErrorMissingConfig        = errors.New("missing VCS configuration")
	ErrorProviderAPI          = errors.New("VCS provider API error")
)
