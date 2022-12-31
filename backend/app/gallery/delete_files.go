package app_gallery

import (
	app_helpers "attiladudas/backend/app/helpers"
	"attiladudas/backend/components/gallery"
	"github.com/gin-gonic/gin"
	"net/http"
)

type deleteFilesBody struct {
	Ids []uint `json:"ids" binding:"required"`
}

func DeleteFilesHandler(fileStore gallery.IFileStore, galleryStore gallery.IGalleryStore) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		deleteFiles(ctx, fileStore, galleryStore)
	}
}

func deleteFiles(ctx *gin.Context, fileStore gallery.IFileStore, galleryStore gallery.IGalleryStore) {
	input := &galleryIdUri{}
	if err := ctx.ShouldBindUri(input); err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	gallery, err := galleryStore.GetGallery(&gallery.GetGalleryInput{Id: &input.Id})
	if app_helpers.HandleError(ctx, err) {
		return
	}

	ids := &deleteFilesBody{}
	if err := ctx.ShouldBindJSON(ids); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	if err := fileStore.DeleteFiles(gallery, ids.Ids); err != nil {
		ctx.Status(http.StatusInternalServerError)
	}
}
