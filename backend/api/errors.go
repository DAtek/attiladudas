package api

import (
	"fibertools"

	"github.com/DAtek/golidator"
)

const (
	TypeErrorRequired         = golidator.ErrorType("REQUIRED")
	TypeErrorAlreadyExists    = golidator.ErrorType("ALREADY_EXISTS")
	TypeErrorInvalid          = golidator.ErrorType("INVALID")
	TypeErrorNotExists        = golidator.ErrorType("NOT_EXISTS")
	TypeErrorWrongCredentials = golidator.ErrorType("WRONG_CREDENTIALS")
)

var (
	ErrorRequired         = fibertools.SimpleValidationError(TypeErrorRequired)
	ErrorWrongCredentials = fibertools.SimpleValidationError(TypeErrorWrongCredentials)
	ErrorInvalid          = fibertools.SimpleValidationError(TypeErrorInvalid)
	ErrorNotExists        = fibertools.SimpleValidationError(TypeErrorNotExists)
	ErrorAlreadyExists    = fibertools.SimpleValidationError(TypeErrorAlreadyExists)
)
