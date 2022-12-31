package app_gallery

import (
	app_helpers "attiladudas/backend/app/helpers"
	"attiladudas/backend/components/gallery"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetGalleryHandler(galleryStore gallery.IGalleryStore, fileStore gallery.IFileStore) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		getGallery(ctx, galleryStore, fileStore)
	}
}

func getGallery(ctx *gin.Context, galleryStore gallery.IGalleryStore, fileStore gallery.IFileStore) {
	input := &getGalleryUriParams{}
	ctx.ShouldBindUri(input)

	getGalleryInput := &gallery.GetGalleryInput{Slug: input.Slug}
	getGalleryInput.SetActive(true)
	gallery, err := galleryStore.GetGallery(getGalleryInput)
	if app_helpers.HandleError(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, convertDbGalleryToApiGallery(gallery, fileStore))
}
