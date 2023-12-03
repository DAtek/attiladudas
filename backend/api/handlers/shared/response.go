package shared

import (
	"api/components/gallery"
	"api/helpers"
	"db/models"
)

type File struct {
	Id       uint   `json:"id"`
	Filename string `json:"filename"`
	Path     string `json:"path"`
	Rank     int    `json:"rank"`
}

type GalleryResponse struct {
	Id          uint    `json:"id"`
	Date        string  `json:"date"`
	Active      bool    `json:"active"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Slug        string  `json:"slug"`
	Files       []*File `json:"files"`
}

func ConvertDbGalleryToApiGallery(gallery *models.Gallery, fileStore gallery.IFileStore) *GalleryResponse {
	files := []*File{}
	for _, file := range gallery.Files {
		files = append(files, &File{
			Id:       file.Id,
			Filename: file.Filename,
			Path:     fileStore.GetFileDownloadPath(gallery, file.Filename),
			Rank:     file.Rank,
		})
	}

	return &GalleryResponse{
		Id:          gallery.Id,
		Title:       gallery.Title,
		Slug:        gallery.Slug,
		Date:        helpers.DateToISO8601(gallery.Date),
		Active:      gallery.Active,
		Description: gallery.Description,
		Files:       files,
	}
}
