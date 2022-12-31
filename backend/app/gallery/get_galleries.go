package app_gallery

import (
	app_helpers "attiladudas/backend/app/helpers"
	"attiladudas/backend/components/auth"
	"attiladudas/backend/components/gallery"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GalleriesResponse struct {
	Galleries []*GalleryResponse `json:"galleries"`
	Total     uint               `json:"total"`
}

type GalleriesQueryParams struct {
	Page     uint `form:"page"`
	PageSize uint `form:"page_size"`
}

func GetGalleriesHandler(
	galleryStore gallery.IGalleryStore,
	fileStore gallery.IFileStore,
	auth auth.IAuthorization,
	jwt auth.IJwt,
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		getGalleries(ctx, galleryStore, fileStore, auth, jwt)
	}
}

func getGalleries(
	ctx *gin.Context,
	galleryStore gallery.IGalleryStore,
	fileStore gallery.IFileStore,
	auth auth.IAuthorization,
	jwt auth.IJwt,
) {
	queryParams := &GalleriesQueryParams{}
	if err := ctx.ShouldBindQuery(&queryParams); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	input := &gallery.GetGalleriesInput{
		Page:     queryParams.Page,
		PageSize: queryParams.PageSize,
	}

	authHeader := ctx.Request.Header.Get("authorization")
	if err := auth.RequireUsername(authHeader, jwt); err != nil {
		input.SetActive(true)
	}

	result, err := galleryStore.GetGalleries(input)

	if app_helpers.HandleError(ctx, err) {
		return
	}

	apiGalleries := []*GalleryResponse{}
	for _, gallery := range result.Items {
		apiGalleries = append(apiGalleries, convertDbGalleryToApiGallery(gallery, fileStore))
	}

	response := &GalleriesResponse{Galleries: apiGalleries, Total: result.Total}
	ctx.JSON(http.StatusOK, response)
}
