package account

import "fmt"

var ErrorInvalidAccountName error
var ErrorReadingStateFile error
var ErrorInvalidDomainValue error

type ErrorInvalidParam struct {
	Param string
}

func (e *ErrorInvalidParam) Error() string {
	return fmt.Sprintf("%s is missing or is not a valid param", e.Param)
}
