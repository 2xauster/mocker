package data

const (
	ErrNotFound = iota 
	ErrAlreadyExists 
	ErrInvalidSyntax
)

// Errors for SQL. So sweet it will give you diabetes.
type SQLError struct {
	Err error
	Code int
}

func NewSQLError(err error, code int) SQLError {
	return SQLError{
		Err: err,
		Code: code,
	}
}

func (err SQLError) Error() string {
	return err.Err.Error()
}