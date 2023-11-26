package fibertools

import (
	"fmt"
	"time"

	"github.com/DAtek/golidator"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type Plugin func(app *fiber.App)

func NewApp(plugins ...Plugin) *fiber.App {
	app := fiber.New(fiber.Config{
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
	})

	app.Use(recover.New())
	app.Use(cors.New())

	loggerConfig := logger.Config{
		Format:     "[${time}] ${method} ${path} ${status} ${latency}\n",
		TimeZone:   "UTC",
		TimeFormat: time.DateTime + ".000",
	}

	app.Use(logger.New(loggerConfig))

	for _, plugin := range plugins {
		plugin(app)
	}

	return app
}
