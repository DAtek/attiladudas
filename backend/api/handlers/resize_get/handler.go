package resize_get

import (
	"api/components/gallery"
	"fibertools"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func PluginResizeImage(mediaDirName string, resizer gallery.IResizer) fibertools.Plugin {
	handler := func(ctx *fiber.Ctx) error {
		return resizeImage(ctx, mediaDirName, resizer)
	}

	return func(app *fiber.App) {
		app.Get(
			fmt.Sprintf("/api/resize/%s/:directory/:size/:filename/", mediaDirName),
			handler,
		)
	}
}

func resizeImage(ctx *fiber.Ctx, mediaDirName string, resizer gallery.IResizer) error {
	params, err := fibertools.BindAndValidateObj[pathData](ctx.ParamsParser)

	if err != nil {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	if _, err := resizer.ResizeImage(params.size, params.Directory, params.Filename); err != nil {
		return err
	}

	return ctx.Status(fiber.StatusFound).Redirect(fmt.Sprintf(
		"/%s/%s/%s/%s",
		mediaDirName,
		params.Directory,
		params.Size,
		params.Filename,
	))
}
