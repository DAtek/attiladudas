package gallery_get

import (
	"api"
	"api/components/gallery"
	"db/models"

	"github.com/DAtek/golidator"
	"github.com/DAtek/gotils"
)

type getGalleryUriParams struct {
	Slug string `params:"slug"`

	gallery *models.Gallery
}

func (obj *getGalleryUriParams) GetValidators(params ...any) golidator.ValidatorCollection {
	galleryStore := params[0].(gallery.IGalleryStore)

	return golidator.ValidatorCollection{
		{Field: "slug", Function: func() *golidator.ValueError {
			getGalleryInput := &gallery.GetGalleryInput{Slug: &obj.Slug}
			getGalleryInput.SetActive(true)
			obj.gallery = gotils.ResultOrPanic(galleryStore.GetGallery(getGalleryInput))

			if obj.gallery == nil {
				return api.ErrorNotExists
			}

			return nil
		}},
	}
}
