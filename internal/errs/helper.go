package errs

import (
	"fmt"
	"strings"
)

func GenericBadRequstErr(what ...string) error {
	whatStr := strings.Join(what, ", ")

	return NewError(fmt.Errorf("%s are/is missing from the request body", whatStr), DataErrorType, ErrDataIllegal)
}