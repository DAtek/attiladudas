package gallery_get

import (
	"api/components/gallery"
	"api/handlers/shared"
	"fibertools"

	"github.com/gofiber/fiber/v2"
)

const path = "/api/gallery/:slug/"

func PluginGetGallery(galleryStore gallery.IGalleryStore, fileStore gallery.IFileStore) fibertools.Plugin {
	var handler = func(ctx *fiber.Ctx) error {
		return getGallery(ctx, galleryStore, fileStore)
	}

	return func(app *fiber.App) {
		app.Get(path, handler)
	}
}

func getGallery(ctx *fiber.Ctx, galleryStore gallery.IGalleryStore, fileStore gallery.IFileStore) error {
	data, err := fibertools.BindAndValidateObj[getGalleryUriParams](ctx.ParamsParser, galleryStore)
	if err != nil {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	respData := shared.ConvertDbGalleryToApiGallery(data.gallery, fileStore)
	return ctx.JSON(respData)
}
