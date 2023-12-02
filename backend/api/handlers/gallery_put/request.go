package gallery_put

import (
	"api"
	"api/components/gallery"
	"db/models"

	"github.com/DAtek/golidator"
	"github.com/DAtek/gotils"
)

type putGalleryPath struct {
	Id uint `params:"id"`

	gallery *models.Gallery
}

func (obj *putGalleryPath) GetValidators(params ...any) golidator.ValidatorCollection {
	galleryStore := params[0].(gallery.IGalleryStore)

	return golidator.ValidatorCollection{
		{Field: "id", Function: func() *golidator.ValueError {
			input := &gallery.GetGalleryInput{Id: &obj.Id}
			obj.gallery = gotils.ResultOrPanic(galleryStore.GetGallery(input))

			if obj.gallery == nil {
				return api.ErrorNotExists
			}

			return nil
		}},
	}
}
