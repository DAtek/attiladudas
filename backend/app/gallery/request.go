package app_gallery

import (
	"attiladudas/backend/components/gallery"
	"attiladudas/backend/helpers"
	"attiladudas/backend/models"

	"github.com/DAtek/golidator"
	"gorm.io/datatypes"
)

type galleryIdUri struct {
	Id uint `uri:"id" binding:"required"`
}

type getGalleryUriParams struct {
	Slug *string `uri:"slug"`
}

type createUpdateGalleryBody struct {
	Title       string `json:"title"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Date        string `json:"date"`
	Active      bool   `json:"active"`
	date        *datatypes.Date
}

func (obj *createUpdateGalleryBody) GetValidators(params ...interface{}) golidator.ValidatorCollection {
	galleryStore, ok := params[0].(gallery.IGalleryStore)
	if !ok {
		panic("GALLERY_STORE_NOT_PROVIDED")
	}

	var originalGallery *models.Gallery = nil
	if len(params) == 2 {
		originalGallery, ok = params[1].(*models.Gallery)
		if !ok {
			panic("ORIGINAL_GALLERY_IS_INVALID")
		}
	}

	return golidator.ValidatorCollection{
		{Field: "title", Function: func() *golidator.ValueError {
			if obj.Title == "" {
				return &golidator.ValueError{ErrorType: helpers.ErrorRequired}
			}

			exists, err := galleryStore.GalleryExists(&gallery.GetGalleryInput{Title: &obj.Title})
			if err != nil {
				panic(err)
			}

			if !exists {
				return nil
			}

			if originalGallery == nil {
				return &golidator.ValueError{ErrorType: helpers.ErrorAlreadyExists}
			}

			if originalGallery != nil && originalGallery.Title != obj.Title {
				return &golidator.ValueError{ErrorType: helpers.ErrorAlreadyExists}
			}

			return nil
		}},
		{Field: "slug", Function: func() *golidator.ValueError {
			if obj.Slug == "" {
				return &golidator.ValueError{ErrorType: helpers.ErrorRequired}
			}

			exists, err := galleryStore.GalleryExists(&gallery.GetGalleryInput{Slug: &obj.Slug})
			if err != nil {
				panic(err)
			}

			if !exists {
				return nil
			}

			if originalGallery == nil {
				return &golidator.ValueError{ErrorType: helpers.ErrorAlreadyExists}
			}

			if originalGallery != nil && originalGallery.Slug != obj.Slug {
				return &golidator.ValueError{ErrorType: helpers.ErrorAlreadyExists}
			}

			return nil
		}},
		{Field: "date", Function: func() *golidator.ValueError {
			if obj.Date == "" {
				return &golidator.ValueError{ErrorType: helpers.ErrorRequired}
			}

			date, err := helpers.DateFromISO8601(obj.Date)
			if err != nil {
				return &golidator.ValueError{ErrorType: helpers.ErrorInvalidISO8601}
			}
			obj.date = date
			return nil
		}},
	}
}
