package gallery_mocks

import (
	"attiladudas/backend/components/gallery"
	"attiladudas/backend/models"
)

type MockGalleryStore struct {
	CreateGallery_ func(*gallery.CreateUpdateGalleryInput) (*models.Gallery, error)
	UpdateGallery_ func(uint, *gallery.CreateUpdateGalleryInput) error
	DeleteGallery_ func(*models.Gallery) error
	GetGallery_    func(input *gallery.GetGalleryInput) (*models.Gallery, error)
	GetGalleries_  func(input *gallery.GetGalleriesInput) (*gallery.PaginatedGalleriesResult, error)
	GalleryExists_ func(input *gallery.GetGalleryInput) (bool, error)
}

func (s *MockGalleryStore) CreateGallery(i *gallery.CreateUpdateGalleryInput) (*models.Gallery, error) {
	return s.CreateGallery_(i)
}

func (s *MockGalleryStore) UpdateGallery(galleryId uint, i *gallery.CreateUpdateGalleryInput) error {
	return s.UpdateGallery_(galleryId, i)
}

func (s *MockGalleryStore) DeleteGallery(gallery *models.Gallery) error {
	return s.DeleteGallery_(gallery)
}

func (s *MockGalleryStore) GetGallery(input *gallery.GetGalleryInput) (*models.Gallery, error) {
	return s.GetGallery_(input)
}

func (s *MockGalleryStore) GetGalleries(input *gallery.GetGalleriesInput) (*gallery.PaginatedGalleriesResult, error) {
	return s.GetGalleries_(input)
}

func (s *MockGalleryStore) GalleryExists(input *gallery.GetGalleryInput) (bool, error) {
	return s.GalleryExists_(input)
}
