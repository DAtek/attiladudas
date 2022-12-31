package helpers

import "github.com/DAtek/golidator"

const ErrorRequired = golidator.ErrorType("REQUIRED")
const ErrorInvalidISO8601 = golidator.ErrorType("INVALID_ISO8601")
const ErrorAlreadyExists = golidator.ErrorType("ALREADY_EXISTS")
const ErrorInvalid = golidator.ErrorType("INVALID")
const ErrorNotExists = golidator.ErrorType("NOT_EXISTS")

type ValidationError struct {
	Location *string              `json:"location"`
	Type     *golidator.ErrorType `json:"type"`
	Context  *map[string]any      `json:"context"`
	Message  *string              `json:"message"`
}

type JsonErrorCollection struct {
	Errors []*ValidationError `json:"errors"`
}

func JsonErrorFromValidationError(err *golidator.ValidationError) *JsonErrorCollection {
	jsonErr := &JsonErrorCollection{}
	errors := []*ValidationError{}

	for _, validationError := range err.Errors {
		errors = append(errors, &ValidationError{
			Type:     &validationError.ErrorType,
			Context:  &validationError.Context,
			Message:  &validationError.Message,
			Location: &validationError.Location,
		})
	}
	jsonErr.Errors = errors
	return jsonErr
}
