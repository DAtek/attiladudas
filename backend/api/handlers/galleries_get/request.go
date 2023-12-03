package galleries_get

import (
	"api"

	"github.com/DAtek/golidator"
)

type GalleriesQueryParams struct {
	Page     *uint `query:"page"`
	PageSize *uint `query:"page_size"`
}

func (obj *GalleriesQueryParams) GetValidators(params ...any) golidator.ValidatorCollection {
	return golidator.ValidatorCollection{
		{Field: "page", Function: func() *golidator.ValueError {
			if obj.Page == nil {
				return api.ErrorInvalid
			}
			return nil
		}},
		{Field: "page_size", Function: func() *golidator.ValueError {
			if obj.PageSize == nil {
				return api.ErrorInvalid
			}
			return nil
		}},
	}
}
