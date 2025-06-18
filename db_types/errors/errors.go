package errors

import "fmt"

type APIErrorCode byte

const (
	ErrCodeMissingParam APIErrorCode = iota + 1
	ErrCodeIncorrectParam
	ErrCodeSQLNoRows
	ErrCodeSQLInternalError
)

type APIError struct {
	Code APIErrorCode `json:"code,omitempty"`
	Err  string       `json:"error,omitempty"`
}

func (e *APIError) Error() string {
	if e == nil {
		return "nil"
	}
	return e.Err
}

var (
	ErrEmptyParam       = NewErrEmptyParam("")
	ErrIncorrectParam   = NewErrIncorrectParam("")
	ErrSQLNoRows        = NewErrSQLNoRows("")
	ErrSQLInternalError = NewErrSQLInternalError("")
)

func NewErrEmptyParam(param string) *APIError {
	return &APIError{
		Code: ErrCodeMissingParam,
		Err:  fmt.Sprintf("missing param: %s", param),
	}
}

func NewErrIncorrectParam(param string) *APIError {
	return &APIError{
		Code: ErrCodeIncorrectParam,
		Err:  fmt.Sprintf("incorrect param: %s", param),
	}
}

func NewErrSQLNoRows(err string) *APIError {
	return &APIError{
		Code: ErrCodeSQLNoRows,
		Err:  err,
	}
}

func NewErrSQLInternalError(err string) *APIError {
	return &APIError{
		Code: ErrCodeSQLInternalError,
		Err:  err,
	}
}

func IsChildErr(err, target error) bool {
	e, ok := err.(*APIError)
	if !ok {
		return false
	}

	t, ok := target.(*APIError)
	if !ok {
		return false
	}

	if e == nil || t == nil {
		return false
	}

	return e.Code == t.Code
}
