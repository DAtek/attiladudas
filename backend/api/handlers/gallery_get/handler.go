package gallery_get

import (
	"api/components/gallery"
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
	data, err := fibertools.BindAndValidateParams(&getGalleryUriParams{}, ctx)
	if err != nil {
		panic(err)
	}

	getGalleryInput := &gallery.GetGalleryInput{Slug: &data.Slug}
	getGalleryInput.SetActive(true)
	gallery, storeErr := galleryStore.GetGallery(getGalleryInput)
	if storeErr != nil {
		panic(err)
	}

	if gallery == nil {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	respData := ConvertDbGalleryToApiGallery(gallery, fileStore)
	return ctx.JSON(respData)
}
