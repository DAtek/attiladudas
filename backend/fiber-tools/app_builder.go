package fibertools

import (
	"fmt"

	"github.com/DAtek/golidator"
	"github.com/gofiber/fiber/v2"
)

type Plugin func(app *fiber.App)

var DefaultConfig = fiber.Config{
	Prefork: false,
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		switch err {
		case ErrorForbidden:
			c.Status(fiber.StatusForbidden)
		case ErrorUnauthorized:
			c.Status(fiber.StatusUnauthorized)
		default:
			validationErr, ok := err.(*golidator.ValidationError)
			if ok {
				return c.Status(fiber.StatusBadRequest).JSON(JsonErrorFromValidationError(validationErr))
			}

			fmt.Printf("Unexpected error: %s\n", err)
			return err
		}
		return nil
	},
}

func NewWithDefaultConfigApp(plugins ...Plugin) *fiber.App {
	return NewApp(DefaultConfig, plugins...)
}

func NewApp(config fiber.Config, plugins ...Plugin) *fiber.App {
	app := fiber.New(config)

	for _, plugin := range plugins {
		plugin(app)
	}

	return app
}
