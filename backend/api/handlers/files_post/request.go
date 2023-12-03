package files_post

import (
	"api"
	"db/models"

	"github.com/DAtek/golidator"
)

type fileInput struct {
	Files []*file
}

func (obj *fileInput) GetValidators(params ...interface{}) golidator.ValidatorCollection {
	return golidator.GetValidatorsForList("files", obj.Files, params...)
}

type file struct {
	Filename string
}

func (obj *file) GetValidators(params ...interface{}) golidator.ValidatorCollection {
	existingFiles := params[0].([]*models.File)
	return []*golidator.Validator{
		{Field: "filename", Function: func() *golidator.ValueError {
			for _, file := range existingFiles {
				if obj.Filename == file.Filename {
					return api.ErrorAlreadyExists
				}
			}
			return nil
		}},
	}
}
