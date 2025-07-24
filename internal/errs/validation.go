package errs

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator"
)

type ValidationErrorResponse struct {
	FailedField string      `json:"failed_field"`
	Tag         string      `json:"tag"`
	Value       interface{} `json:"value"`
}

func (err ValidationErrorResponse) Error() string {
	return fmt.Sprintf("validation error : [FailedField : %s] [Tag : %s] [Value : %v]", err.FailedField, err.Tag, err.Value)
}

var validate = validator.New()

func Validate(data interface{}) (error) {
	validationErrs := []ValidationErrorResponse{}

	vErrs := validate.Struct(data)
	if vErrs != nil {
		for _, err := range vErrs.(validator.ValidationErrors) {
			var elem ValidationErrorResponse

			elem.FailedField = strings.ToLower(err.Field())
			elem.Tag = strings.ToLower(err.Tag())
			elem.Value = err.Value()

			validationErrs = append(validationErrs, elem)
		}
	}

	if len(validationErrs) == 0 {
		return nil
	}

	errs := make([]error, 0, len(validationErrs))
	for _, ve := range validationErrs {
		errs = append(errs, ve)
	}

	newE := NewError(errors.Join(errs...), DataErrorType, ErrDataIllegal)
	return newE
}
