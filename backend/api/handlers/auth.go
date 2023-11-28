package handlers

import (
	"api/components/auth"
	"fibertools"

	"github.com/gofiber/fiber/v2"
)

func RequireUsername(authCtx auth.IAuthorization) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if err := authCtx.RequireUsername(authHeader); err != nil {
			return fibertools.ErrorUnauthorized
		}
		return c.Next()
	}
}
