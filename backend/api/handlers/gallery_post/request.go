package gallery_post

import (
	"api"
	"api/components/gallery"
	"api/helpers"
	"db/models"

	"github.com/DAtek/golidator"
	"gorm.io/datatypes"
)

type createUpdateGalleryBody struct {
	Title       string `json:"title"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Date        string `json:"date"`
	Active      bool   `json:"active"`
	date        *datatypes.Date
}

func (obj *createUpdateGalleryBody) GetValidators(params ...interface{}) golidator.ValidatorCollection {
	galleryStore := params[0].(gallery.IGalleryStore)
	var originalGallery *models.Gallery = nil

	if len(params) == 2 {
		originalGallery = params[1].(*models.Gallery)
	}
	return golidator.ValidatorCollection{
		{Field: "title", Function: func() *golidator.ValueError {
			if obj.Title == "" {
				return api.ErrorRequired
			}

			exists, err := galleryStore.GalleryExists(&gallery.GetGalleryInput{Title: &obj.Title})

			if err != nil {
				panic(err)
			}

			if !exists {
				return nil
			}

			if originalGallery == nil {
				return api.ErrorAlreadyExists
			}

			if originalGallery != nil && originalGallery.Title != obj.Title {
				return api.ErrorAlreadyExists
			}

			return nil
		}},
		{Field: "slug", Function: func() *golidator.ValueError {
			if obj.Slug == "" {
				return api.ErrorRequired
			}

			exists, err := galleryStore.GalleryExists(&gallery.GetGalleryInput{Slug: &obj.Slug})
			if err != nil {
				panic(err)
			}

			if !exists {
				return nil
			}

			if originalGallery == nil {
				return api.ErrorAlreadyExists
			}

			if originalGallery != nil && originalGallery.Slug != obj.Slug {
				return api.ErrorAlreadyExists
			}

			return nil
		}},
		{Field: "date", Function: func() *golidator.ValueError {
			if obj.Date == "" {
				return api.ErrorRequired
			}

			date, err := helpers.DateFromISO8601(obj.Date)
			if err != nil {
				return api.ErrorInvalid
			}
			obj.date = date

			return nil
		}},
	}
}
