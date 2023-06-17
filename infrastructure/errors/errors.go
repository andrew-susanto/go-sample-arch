package errors

import "errors"

// Error wraps golang error with error type and error code
type Error struct {
	E       error
	EType   int
	ECode   string
	ENumber int // for JSON RPC error number ex. -32000
}

const (
	USER   int = 1
	SYSTEM int = 2
)

// Wrap converts golang standard error type to custom error type
func Wrap(e error) Error {
	switch val := e.(type) {
	case Error:
		return val
	}

	if e == nil {
		return Wrap(New("unknown error").WithCode("INF.ERR00"))
	}

	return Error{
		E:     e,
		EType: SYSTEM, // default error system
	}
}

// New returns an error that formats as the given text. Each call to New returns a distinct error value even if the text is identical.
func New(msg string) Error {
	return Error{
		E:     errors.New(msg),
		EType: SYSTEM, // default error system
	}
}

// WithType update error with given type
func (err Error) WithType(eType int) Error {
	err.EType = eType
	return err
}

// WithCode update error with given code
func (err Error) WithCode(eCode string) Error {
	if err.ECode != "" {
		return err
	}

	err.ECode = eCode
	return err
}

// WithNumber update error with given number
func (err Error) WithNumber(eCode int) Error {
	err.ENumber = eCode
	return err
}

// Error returns error string
func (err Error) Error() string {
	return err.E.Error()
}

// Is reports whether any error in err's tree matches target.
func Is(e error, target error) bool {
	return errors.Is(e, target)
}
