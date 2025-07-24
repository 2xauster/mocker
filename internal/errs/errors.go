package errs

import "encoding/json"

const (
	ErrNotFound = iota
	ErrAlreadyExists
	ErrInvalidSyntax

	ErrDataMismatch // (x1, x2) -> x1 != x2
	ErrDataIllegal  // Data that violates the schema.
	ErrInternalFailure
)

type ErrorType int

const (
	SQLErrorType  ErrorType = iota // This is for SQL related errors.
	DataErrorType                  // DataErrorType is a general error type.
)

func (et ErrorType) String() string {
	return [...]string{"SQLErrorType", "DataErrorType"}[et]
}

type Error struct {
	Err  error  `json:"-"`
	Type string `json:"type"`
	Code int    `json:"code"`
}

func NewError(err error, errType ErrorType, code int) Error {
	return Error{
		Err:  err,
		Type: errType.String(),
		Code: code,
	}
}

func (err Error) Error() string {
	return err.Err.Error()
}

func (err Error) Unwrap() error {
	return err.Err
}

func (err Error) Is(target error) bool {
	t, ok := target.(Error)
	if !ok {
		return false
	}
	return err.Code == t.Code && err.Type == t.Type
}

func (err Error) MarshalJSON() ([]byte, error) {
	type Alias Error
	return json.Marshal(&struct {
		Alias
		Details string `json:"details"`
	}{
		Alias:   (Alias)(err),
		Details: err.Err.Error(),
	})
}