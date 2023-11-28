package fibertools

import (
	"github.com/DAtek/golidator"
	"github.com/gofiber/fiber/v2"
)

func BindAndValidateJSON[T golidator.IValidators](obj T, ctx *fiber.Ctx, validatorContext ...any) (T, *golidator.ValidationError) {
	ctx.BodyParser(obj)
	return validateObj(obj, validatorContext...)
}

func BindAndValidateQuery[T golidator.IValidators](obj T, ctx *fiber.Ctx, validatorContext ...any) (T, *golidator.ValidationError) {
	ctx.QueryParser(obj)
	return validateObj(obj, validatorContext...)
}

func BindAndValidateParams[T golidator.IValidators](obj T, ctx *fiber.Ctx, validatorContext ...any) (T, *golidator.ValidationError) {
	ctx.ParamsParser(obj)
	return validateObj(obj, validatorContext...)
}

func validateObj[T golidator.IValidators](obj T, validatorContext ...any) (T, *golidator.ValidationError) {
	if err := golidator.Validate(obj, validatorContext...); err != nil {
		return obj, err
	}

	return obj, nil
}
