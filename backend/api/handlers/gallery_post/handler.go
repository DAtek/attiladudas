package gallery_post

import (
	"api/components/auth"
	"api/components/gallery"
	"api/handlers"
	"fibertools"

	"github.com/DAtek/gotils"
	"github.com/gofiber/fiber/v2"
)

const path = "/api/gallery/"

type createGalleryResponse struct {
	Id uint `json:"id"`
}

func PluginPostGallery(authCtx auth.IAuthorization, galleryStore gallery.IGalleryStore) fibertools.Plugin {
	handler := func(ctx *fiber.Ctx) error {
		return postGallery(ctx, galleryStore)
	}

	return func(app *fiber.App) {
		app.Post(path, handlers.RequireUsername(authCtx), handler)
	}
}

func postGallery(ctx *fiber.Ctx, galleryStore gallery.IGalleryStore) error {
	requestBody, err := fibertools.BindAndValidateJSON(&createUpdateGalleryBody{}, ctx, galleryStore)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fibertools.JsonErrorFromValidationError(err))
	}

	gallery := gotils.ResultOrPanic(galleryStore.CreateGallery(&gallery.CreateUpdateGalleryInput{
		Title:       requestBody.Title,
		Slug:        requestBody.Slug,
		Description: requestBody.Description,
		Date:        requestBody.date,
		Active:      requestBody.Active,
	}))

	return ctx.Status(fiber.StatusCreated).JSON(&createGalleryResponse{Id: gallery.Id})
}
