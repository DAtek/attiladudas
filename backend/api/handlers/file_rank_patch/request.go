package file_rank_patch

import (
	"api"
	"api/components/gallery"

	"github.com/DAtek/golidator"
)

type patchFileRankPathParams struct {
	FileId *uint `params:"id"`
	Rank   *int  `params:"rank"`
}

func (obj *patchFileRankPathParams) GetValidators(ctx ...interface{}) golidator.ValidatorCollection {
	fileStore := ctx[0].(gallery.IFileStore)

	return []*golidator.Validator{
		{Field: "id", Function: func() *golidator.ValueError {
			file, err := fileStore.GetFile(*obj.FileId)
			if err != nil {
				panic(err)
			}
			if file == nil {
				return api.ErrorNotExists
			}

			return nil
		}},
		{Field: "rank", Function: func() *golidator.ValueError {
			return nil
		}},
	}
}
