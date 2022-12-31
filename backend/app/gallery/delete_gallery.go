package app_gallery

import (
	app_helpers "attiladudas/backend/app/helpers"
	"attiladudas/backend/components/gallery"
	"net/http"

	"github.com/gin-gonic/gin"
)

type deleteGalleryInput struct {
	Id uint `uri:"id" binding:"required"`
}

func DeleteGalleryHandler(store gallery.IGalleryStore) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		deleteGallery(ctx, store)
	}
}

func deleteGallery(ctx *gin.Context, store gallery.IGalleryStore) {
	input := &deleteGalleryInput{}
	if err := ctx.ShouldBindUri(input); err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	gallery, err := store.GetGallery(&gallery.GetGalleryInput{Id: &input.Id})
	if app_helpers.HandleError(ctx, err) {
		return
	}

	if err := store.DeleteGallery(gallery); err != nil {
		ctx.Status(http.StatusInternalServerError)
	}
}
