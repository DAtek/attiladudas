package gallery_delete

import (
	"api/components/auth"
	"api/components/gallery"
	"api/handlers"
	"api/handlers/shared"
	"fibertools"

	"github.com/gofiber/fiber/v2"
)

func PluginDeleteGallery(authCtx auth.IAuthorization, galleryStore gallery.IGalleryStore) fibertools.Plugin {
	handler := func(ctx *fiber.Ctx) error {
		return deleteGallery(ctx, galleryStore)
	}

	return func(app *fiber.App) {
		app.Delete("/api/gallery/:id/", handlers.RequireUsername(authCtx), handler)
	}
}

func deleteGallery(ctx *fiber.Ctx, galleryStore gallery.IGalleryStore) error {
	pathData, err := fibertools.BindAndValidateObj[shared.GalleryIdInPath](ctx.ParamsParser, galleryStore)

	if err != nil {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	return galleryStore.DeleteGallery(pathData.Gallery)
}
