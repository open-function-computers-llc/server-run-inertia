package errx

import "fmt"

// Wrapf returns err as a new error that preserves the original error type
// (so errors.Is/As still match) with a human-readable verb and cause appended, e.g.:
//
//	errx.Wrapf("reading settings file", err, ErrorReadingStateFile)
//	-> "reading settings file: open /path/settings.json: no such file or directory"
func Wrapf(verb string, err, sentinel error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w: %w", verb, sentinel, err)
}
