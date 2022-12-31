package app_helpers

import (
	"attiladudas/backend/components"
	"net/http"

	"github.com/DAtek/golidator"
	"github.com/gin-gonic/gin"
)

func HandleError(ctx *gin.Context, err error) bool {
	switch err {
	case nil:
		return false
	case components.NotFoundError:
		ctx.Status(http.StatusNotFound)
		return true
	default:
		ctx.Status(http.StatusInternalServerError)
		return true
	}
}

func BindAndValidateJSON[T golidator.IValidators](obj T, ctx *gin.Context, validatorContext ...interface{}) (T, *golidator.ValidationError) {
	ctx.ShouldBindJSON(obj)
	if err := golidator.Validate(obj, validatorContext...); err != nil {
		return obj, err
	}

	return obj, nil
}
