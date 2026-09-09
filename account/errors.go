package account

import "errors"

// sentinel errors for use with errors.Is
var (
	ErrorInvalidAccountName = errors.New("invalid account name")
	ErrorReadingStateFile   = errors.New("error reading account state file")
	ErrorInvalidDomainValue = errors.New("invalid domain value")
	ErrorInvalidParam       = errors.New("invalid param provided")
)
