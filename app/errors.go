package app

type ErrorInvalidArguments struct{}

func (e *ErrorInvalidArguments) Error() string {
	return "Sorry, but you didn't provide valid arguments"
}
