package api

import (
	"fibertools"

	"github.com/DAtek/golidator"
)

const (
	errorRequired         = golidator.ErrorType("REQUIRED")
	errorInvalidISO8601   = golidator.ErrorType("INVALID_ISO8601")
	errorAlreadyExists    = golidator.ErrorType("ALREADY_EXISTS")
	errorInvalid          = golidator.ErrorType("INVALID")
	errorNotExists        = golidator.ErrorType("NOT_EXISTS")
	errorWrongCredentials = golidator.ErrorType("WRONG_CREDENTIALS")
)

var (
	ErrorRequired         = fibertools.SimpleValidationError(errorRequired)
	ErrorWrongCredentials = fibertools.SimpleValidationError(errorWrongCredentials)
	ErrorInvalid          = fibertools.SimpleValidationError(errorInvalid)
	ErrorNotExists        = fibertools.SimpleValidationError(errorNotExists)
)
