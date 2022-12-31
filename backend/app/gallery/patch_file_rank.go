package app_gallery

import (
	"attiladudas/backend/components/gallery"
	"attiladudas/backend/helpers"
	"net/http"

	"github.com/DAtek/golidator"
	"github.com/gin-gonic/gin"
)

type patchFileRankPathParams struct {
	FileId *uint `uri:"id"`
	Rank   *int  `uri:"rank"`
}

func (obj *patchFileRankPathParams) GetValidators(ctx ...interface{}) golidator.ValidatorCollection {
	fileStore, ok := ctx[0].(gallery.IFileStore)
	if !ok {
		panic("File store not provided for patch file rank validation")
	}
	return []*golidator.Validator{
		{Field: "id", Function: func() *golidator.ValueError {
			if obj.FileId == nil {
				return &golidator.ValueError{ErrorType: helpers.ErrorInvalid}
			}
			file, err := fileStore.GetFile(*obj.FileId)
			if err != nil {
				panic(err)
			}
			if file == nil {
				return &golidator.ValueError{ErrorType: helpers.ErrorNotExists}
			}

			return nil
		}},
		{Field: "rank", Function: func() *golidator.ValueError {
			if obj.Rank == nil {
				return &golidator.ValueError{ErrorType: helpers.ErrorInvalid}
			}
			return nil
		}},
	}
}

func PatchFileRankHandler(store gallery.IFileStore) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		patchFileRank(ctx, store)
	}
}

func patchFileRank(ctx *gin.Context, store gallery.IFileStore) {
	input := &patchFileRankPathParams{}
	ctx.ShouldBindUri(input)
	if err := golidator.Validate(input, store); err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	err := store.UpdateFileRank(&gallery.UpdateFileRankInput{FileId: *input.FileId, Rank: *input.Rank})
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
	}
}
