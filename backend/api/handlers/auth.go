package handlers

import (
	"api/components/auth"

	"github.com/gofiber/fiber/v2"
)

func CreateAuthHandler(authCtx auth.IAuthorization) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		authHeader := ctx.Get("Authorization")
		if err := authCtx.RequireUsername(authHeader); err != nil {
			return ctx.SendStatus(fiber.StatusUnauthorized)
		}
		return ctx.Next()
	}
}
