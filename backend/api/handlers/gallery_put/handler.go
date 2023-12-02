package gallery_put

import (
	"api/components/auth"
	"api/components/gallery"
	"api/handlers"
	"api/handlers/gallery_post"
	"api/handlers/shared"
	"fibertools"

	"github.com/gofiber/fiber/v2"
)

func PluginPutGallery(authCtx auth.IAuthorization, galleryStore gallery.IGalleryStore) fibertools.Plugin {
	handler := func(ctx *fiber.Ctx) error {
		return putGallery(ctx, galleryStore)
	}
	return func(app *fiber.App) {
		app.Put("/api/gallery/:id/", handlers.RequireUsername(authCtx), handler)
	}
}

func putGallery(ctx *fiber.Ctx, galleryStore gallery.IGalleryStore) error {
	pathParam, err := fibertools.BindAndValidateObj[shared.GalleryIdInPath](ctx.ParamsParser, galleryStore)
	if err != nil {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	requestBody, validationErr := fibertools.BindAndValidateObj[gallery_post.CreateUpdateGalleryBody](ctx.BodyParser, galleryStore, pathParam.Gallery)
	if validationErr != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fibertools.JsonErrorFromValidationError(validationErr))
	}

	if err := galleryStore.UpdateGallery(pathParam.Id, &gallery.CreateUpdateGalleryInput{
		Title:       requestBody.Title,
		Slug:        requestBody.Slug,
		Description: requestBody.Description,
		Date:        requestBody.ParsedDate,
		Active:      requestBody.Active,
	}); err != nil {
		panic(err)
	}

	return nil
}
