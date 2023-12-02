package fibertools

import (
	"github.com/gofiber/fiber/v2"
)

type Plugin func(app *fiber.App)

var DefaultConfig = fiber.Config{}

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
