package files_delete

import (
	"api/components/auth"
	"api/components/gallery"
	"api/handlers"
	"api/handlers/shared"
	"fibertools"

	"github.com/gofiber/fiber/v2"
)

func PluginDeleteFiles(authCtx auth.IAuthorization, galleryStore gallery.IGalleryStore, fileStore gallery.IFileStore) fibertools.Plugin {
	handler := func(ctx *fiber.Ctx) error {
		return deleteFiles(ctx, galleryStore, fileStore)
	}

	return func(app *fiber.App) {
		app.Delete("/api/gallery/:id/files/", handlers.CreateAuthHandler(authCtx), handler)
	}
}

func deleteFiles(ctx *fiber.Ctx, galleryStore gallery.IGalleryStore, fileStore gallery.IFileStore) error {
	pathParam, err := fibertools.BindAndValidateObj[shared.GalleryIdInPath](ctx.ParamsParser, galleryStore)
	if err != nil {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	ids, err := fibertools.BindAndValidateObj[deleteFilesBody](ctx.BodyParser)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fibertools.JsonErrorFromValidationError(err))
	}

	if err := fileStore.DeleteFiles(pathParam.Gallery, ids.Ids); err != nil {
		panic(err)
	}

	return nil
}
