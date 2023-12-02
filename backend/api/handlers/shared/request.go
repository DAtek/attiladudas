package shared

import (
	"api"
	"api/components/gallery"
	"db/models"

	"github.com/DAtek/golidator"
	"github.com/DAtek/gotils"
)

type GalleryIdInPath struct {
	Id uint `params:"id"`

	Gallery *models.Gallery
}

func (obj *GalleryIdInPath) GetValidators(params ...any) golidator.ValidatorCollection {
	galleryStore := params[0].(gallery.IGalleryStore)

	return golidator.ValidatorCollection{
		{Field: "id", Function: func() *golidator.ValueError {
			input := &gallery.GetGalleryInput{Id: &obj.Id}
			obj.Gallery = gotils.ResultOrPanic(galleryStore.GetGallery(input))

			if obj.Gallery == nil {
				return api.ErrorNotExists
			}

			return nil
		}},
	}
}
