package xerr

import (
	"errors"
)

// StackTraceError decorates an error with a stack trace.
type StackTraceError struct {
	root       error
	stackTrace StackTrace
}

func (e *StackTraceError) Error() string {
	return e.root.Error()
}

func (e *StackTraceError) StackTrace() StackTrace {
	return e.stackTrace
}

// GetStackTrace returns the stack trace of an error. If the error is not a StackTraceError, this function attempts to
// unwrap the error to find a StackTraceError.
func GetStackTrace(err error) StackTrace {
	var ste *StackTraceError

	if errors.As(err, &ste) {
		return ste.stackTrace
	}

	return nil
}

func StackTraceErr(msg string) error {
	return &StackTraceError{
		root:       errors.New(msg),
		stackTrace: NewStackTrace(1),
	}
}

func AddStackTrace(err error) error {
	return &StackTraceError{
		root:       err,
		stackTrace: NewStackTrace(1),
	}
}
