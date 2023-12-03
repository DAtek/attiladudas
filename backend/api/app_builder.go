package api

import (
	"fibertools"

	"github.com/gofiber/fiber/v2"
)

func AppWithMiddlewares(plugins ...fibertools.Plugin) *fiber.App {
	allPlugins := []fibertools.Plugin{PluginMiddlewares()}
	allPlugins = append(allPlugins, plugins...)
	config := fiber.Config{
		BodyLimit: 50 * 1024 * 1024,
	}
	return fibertools.NewApp(config, allPlugins...)
}
