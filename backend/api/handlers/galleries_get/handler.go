package galleries_get

import (
	"api/components/auth"
	"api/components/gallery"
	"api/handlers/gallery_get"
	"fibertools"

	"github.com/DAtek/gotils"
	"github.com/gofiber/fiber/v2"
)

const path = "/api/galleries/"

type GalleriesResponse struct {
	Galleries []*gallery_get.GalleryResponse `json:"galleries"`
	Total     uint                           `json:"total"`
}

func PluginGetGalleries(
	galleryStore gallery.IGalleryStore,
	fileStore gallery.IFileStore,
	authCtx auth.IAuthorization,
) fibertools.Plugin {
	handler := func(ctx *fiber.Ctx) error {
		return getGalleries(ctx, galleryStore, fileStore, authCtx)
	}
	return func(app *fiber.App) {
		app.Get(path, handler)
	}
}

func getGalleries(
	ctx *fiber.Ctx,
	galleryStore gallery.IGalleryStore,
	fileStore gallery.IFileStore,
	authCtx auth.IAuthorization,
) error {
	params, err := fibertools.BindAndValidateObj[GalleriesQueryParams](ctx.QueryParser)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fibertools.JsonErrorFromValidationError(err))
	}

	active := false
	if err := authCtx.RequireUsername(ctx.Get("Authorization")); err != nil {
		active = true
	}

	result := gotils.ResultOrPanic(galleryStore.GetGalleries(&gallery.GetGalleriesInput{
		Active:   &active,
		Page:     *params.Page,
		PageSize: *params.PageSize,
	}))

	apiGalleries := []*gallery_get.GalleryResponse{}
	for _, gallery := range result.Items {
		apiGalleries = append(apiGalleries, gallery_get.ConvertDbGalleryToApiGallery(gallery, fileStore))
	}

	response := &GalleriesResponse{Galleries: apiGalleries, Total: result.Total}
	return ctx.JSON(response)
}
