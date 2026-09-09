package app

import "errors"

var (
	ErrorInvalidArguments          = errors.New("no valid arguments provided")
	ErrorInvalidSetting            = errors.New("invalid setting requested")
	ErrorMissingRequiredFlag       = errors.New("missing required flag")
	ErrorInvalidParam              = errors.New("invalid param provided")
	ErrorDatabaseNotFoundInAccount = errors.New("database not found in account")
)
