package components

import (
	"fibertools"

	"github.com/DAtek/golidator"
)

const errorRequired = golidator.ErrorType("REQUIRED")
const errorInvalidISO8601 = golidator.ErrorType("INVALID_ISO8601")
const errorAlreadyExists = golidator.ErrorType("ALREADY_EXISTS")
const errorInvalid = golidator.ErrorType("INVALID")
const errorNotExists = golidator.ErrorType("NOT_EXISTS")
const errorWrongCredentials = golidator.ErrorType("WRONG_CREDENTIALS")

type ApiError string

func (e ApiError) Error() string {
	return string(e)
}

var ErrorRequired = fibertools.SimpleValidationError(errorRequired)
var ErrorWrongCredentials = fibertools.SimpleValidationError(errorWrongCredentials)
