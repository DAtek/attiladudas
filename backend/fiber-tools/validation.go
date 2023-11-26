package fibertools

import (
	"github.com/DAtek/golidator"
	"github.com/gofiber/fiber/v2"
)

func BindAndValidateJSON[T golidator.IValidators](obj T, ctx *fiber.Ctx, validatorContext ...any) (T, *golidator.ValidationError) {
	ctx.BodyParser(obj)
	if err := golidator.Validate(obj, validatorContext...); err != nil {
		return obj, err
	}

	return obj, nil
}

func BindAndValidateQuery[T golidator.IValidators](obj T, ctx *fiber.Ctx, validatorContext ...any) (T, *golidator.ValidationError) {
	ctx.QueryParser(obj)
	if err := golidator.Validate(obj, validatorContext...); err != nil {
		return obj, err
	}

	return obj, nil
}
