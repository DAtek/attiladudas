package app_gallery

import (
	app_helpers "attiladudas/backend/app/helpers"
	"attiladudas/backend/components/gallery"
	"attiladudas/backend/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createGalleryResponse struct {
	Id uint `json:"id"`
}

func PostGalleryHandler(galleryStore gallery.IGalleryStore) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		postGallery(ctx, galleryStore)
	}
}

func postGallery(ctx *gin.Context, galleryStore gallery.IGalleryStore) {
	requestBody, err := app_helpers.BindAndValidateJSON(&createUpdateGalleryBody{}, ctx, galleryStore)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.JsonErrorFromValidationError(err))
		return
	}

	gallery, storeErr := galleryStore.CreateGallery(&gallery.CreateUpdateGalleryInput{
		Title:       requestBody.Title,
		Slug:        requestBody.Slug,
		Description: requestBody.Description,
		Date:        requestBody.date,
		Active:      requestBody.Active,
	})

	if storeErr != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusCreated, &createGalleryResponse{Id: gallery.Id})
}
