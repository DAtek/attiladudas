package shared

import (
	"api"
	"api/components/gallery"
	"api/helpers"
	"db/models"

	"github.com/DAtek/golidator"
	"github.com/DAtek/gotils"
	"gorm.io/datatypes"
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

type CreateUpdateGalleryBody struct {
	Title       string `json:"title"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Date        string `json:"date"`
	Active      bool   `json:"active"`

	ParsedDate *datatypes.Date
}

func (obj *CreateUpdateGalleryBody) GetValidators(params ...interface{}) golidator.ValidatorCollection {
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

			exists := gotils.ResultOrPanic(galleryStore.GalleryExists(&gallery.GetGalleryInput{Title: &obj.Title}))

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

			exists := gotils.ResultOrPanic(galleryStore.GalleryExists(&gallery.GetGalleryInput{Slug: &obj.Slug}))

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
			obj.ParsedDate = date

			return nil
		}},
	}
}
