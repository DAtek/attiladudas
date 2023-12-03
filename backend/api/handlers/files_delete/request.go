package files_delete

import "github.com/DAtek/golidator"

type deleteFilesBody struct {
	Ids []uint `json:"ids"`
}

func (obj *deleteFilesBody) GetValidators(params ...any) golidator.ValidatorCollection {
	return golidator.ValidatorCollection{
		{Field: "ids", Function: func() *golidator.ValueError {
			return nil
		}},
	}
}
