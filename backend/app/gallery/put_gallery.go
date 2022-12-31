package app_gallery

import (
	app_helpers "attiladudas/backend/app/helpers"
	"attiladudas/backend/components/gallery"
	"attiladudas/backend/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PutGalleryHandler(store gallery.IGalleryStore) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		putGallery(ctx, store)
	}
}

func putGallery(ctx *gin.Context, store gallery.IGalleryStore) {
	input := &galleryIdUri{}
	if err := ctx.ShouldBindUri(input); err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	galleryObj, err := store.GetGallery(&gallery.GetGalleryInput{Id: &input.Id})
	if app_helpers.HandleError(ctx, err) {
		return
	}

	requestBody, validationErr := app_helpers.BindAndValidateJSON(&createUpdateGalleryBody{}, ctx, store, galleryObj)
	if validationErr != nil {
		ctx.JSON(http.StatusBadRequest, helpers.JsonErrorFromValidationError(validationErr))
		return
	}

	if err := store.UpdateGallery(galleryObj.Id, &gallery.CreateUpdateGalleryInput{
		Title:       requestBody.Title,
		Slug:        requestBody.Slug,
		Description: requestBody.Description,
		Date:        requestBody.date,
		Active:      requestBody.Active,
	}); err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
}
