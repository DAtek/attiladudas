package api

import (
	"fibertools"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func PluginMiddlewares() fibertools.Plugin {
	return func(app *fiber.App) {
		app.Use(recover.New())
		app.Use(cors.New())

		loggerConfig := logger.Config{
			Format:     "[${time}] ${method} ${path} ${status} ${latency}\n",
			TimeZone:   "UTC",
			TimeFormat: time.DateTime + ".000",
		}

		app.Use(logger.New(loggerConfig))
	}
}
