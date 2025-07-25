package errors

type CustomError struct {
	Message string
	Code    int
	Err     error
}

func (e *CustomError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func NewCustomError(message string, code int, err error) *CustomError {
	return &CustomError{Message: message, Code: code, Err: err}
}