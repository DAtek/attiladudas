package api

import (
	"fibertools"

	"github.com/gofiber/fiber/v2"
)

func AppWithMiddlewares(plugins ...fibertools.Plugin) *fiber.App {
	allPlugins := []fibertools.Plugin{PluginMiddlewares()}
	allPlugins = append(allPlugins, plugins...)

	return fibertools.NewWithDefaultConfigApp(allPlugins...)
}
