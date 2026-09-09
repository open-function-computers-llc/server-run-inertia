package account

import (
	"errors"
	"fmt"
)

// sentinel errors for use with errors.Is
var (
	ErrorInvalidAccountName = errors.New("invalid account name")
	ErrorReadingStateFile   = errors.New("error reading account state file")
	ErrorInvalidDomainValue = errors.New("invalid domain value")
)

type ErrorInvalidParam struct {
	Param string
}

func (e *ErrorInvalidParam) Error() string {
	return fmt.Sprintf("%s is missing or is not a valid param", e.Param)
}
