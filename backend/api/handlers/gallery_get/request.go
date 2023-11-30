package gallery_get

import "github.com/DAtek/golidator"

type getGalleryUriParams struct {
	Slug string `params:"slug"`
}

func (obj *getGalleryUriParams) GetValidators(params ...any) golidator.ValidatorCollection {
	return golidator.ValidatorCollection{}
}
