package scriptrunner

import "errors"

var (
	ErrorMissingScript   = errors.New("missing script")
	ErrorInvalidScript   = errors.New("invalid script")
	ErrorMissingArgument = errors.New("missing required script argument")
)
