package fibertools

import "github.com/DAtek/golidator"

type ValidationError struct {
	Location string              `json:"location"`
	Type     golidator.ErrorType `json:"type"`
	Context  map[string]any      `json:"context"`
	Message  string              `json:"message"`
}

type JsonErrorCollection struct {
	Errors []*ValidationError `json:"errors"`
}

func JsonErrorFromValidationError(err *golidator.ValidationError) *JsonErrorCollection {
	jsonErr := &JsonErrorCollection{}
	errors := []*ValidationError{}

	for _, validationError := range err.Errors {
		errors = append(errors, &ValidationError{
			Type:     validationError.ErrorType,
			Context:  validationError.Context,
			Message:  validationError.Message,
			Location: validationError.Location,
		})
	}
	jsonErr.Errors = errors
	return jsonErr
}

func SimpleValidationError(errorType golidator.ErrorType) *golidator.ValueError {
	return &golidator.ValueError{ErrorType: errorType}
}

type AppError string

func (err AppError) Error() string {
	return string(err)
}

const ErrorForbidden = AppError("FORBIDDEN")
const ErrorUnauthorized = AppError("UNAUTHORIZED")
