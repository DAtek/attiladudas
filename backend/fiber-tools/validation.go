package fibertools

import (
	"github.com/DAtek/golidator"
)

const TypeError = golidator.ErrorType("TYPE_ERROR")

func BindAndValidateObj[T any](parseObj parser, validatorContext ...any) (*T, *golidator.ValidationError) {
	var obj T
	if err := parseObj(&obj); err != nil {
		fieldError := &golidator.FieldError{}
		fieldError.Location = "__root__"
		fieldError.ErrorType = TypeError
		fieldError.Message = err.Error()
		return nil, &golidator.ValidationError{Errors: []*golidator.FieldError{fieldError}}
	}

	validable := anyToValidable(&obj)
	if err := golidator.Validate(validable, validatorContext...); err != nil {
		return &obj, err
	}
	return &obj, nil
}

type parser func(obj any) error

func anyToValidable(obj any) golidator.IValidators {
	validable, ok := obj.(golidator.IValidators)
	if !ok {
		panic("FAILED CONVERTING TO golidator.IValidators")
	}
	return validable
}
