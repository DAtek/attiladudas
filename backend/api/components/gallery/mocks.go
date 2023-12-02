package gallery

import "db/models"

type MockFileStore struct {
	AddFiles_            func(*models.Gallery, []*FileInput) error
	DeleteFiles_         func(*models.Gallery, []uint) error
	GetFileDownloadPath_ func(gallery *models.Gallery, filename string) string
	UpdateFileRank_      func(data *UpdateFileRankInput) error
	GetFile_             func(id uint) (*models.File, error)
	MediaDirName_        func() string
}

func (s *MockFileStore) AddFiles(gallery *models.Gallery, files []*FileInput) error {
	return s.AddFiles_(gallery, files)
}

func (s *MockFileStore) DeleteFiles(gallery *models.Gallery, fileIds []uint) error {
	return s.DeleteFiles_(gallery, fileIds)
}

func (s *MockFileStore) GetFileDownloadPath(gallery *models.Gallery, filename string) string {
	return s.GetFileDownloadPath_(gallery, filename)
}

func (s *MockFileStore) UpdateFileRank(data *UpdateFileRankInput) error {
	return s.UpdateFileRank_(data)
}

func (s *MockFileStore) GetFile(id uint) (*models.File, error) {
	return s.GetFile_(id)
}

func (s *MockFileStore) MediaDirName() string {
	return s.MediaDirName_()
}

type MockGalleryStore struct {
	CreateGallery_ func(input *CreateUpdateGalleryInput) (*models.Gallery, error)
	UpdateGallery_ func(galleryId uint, input *CreateUpdateGalleryInput) error
	DeleteGallery_ func(gallery *models.Gallery) error
	GetGallery_    func(input *GetGalleryInput) (*models.Gallery, error)
	GetGalleries_  func(input *GetGalleriesInput) (*PaginatedGalleriesResult, error)
	GalleryExists_ func(input *GetGalleryInput) (bool, error)
}

func (s *MockGalleryStore) CreateGallery(i *CreateUpdateGalleryInput) (*models.Gallery, error) {
	return s.CreateGallery_(i)
}

func (s *MockGalleryStore) UpdateGallery(galleryId uint, i *CreateUpdateGalleryInput) error {
	return s.UpdateGallery_(galleryId, i)
}

func (s *MockGalleryStore) DeleteGallery(gallery *models.Gallery) error {
	return s.DeleteGallery_(gallery)
}

func (s *MockGalleryStore) GetGallery(input *GetGalleryInput) (*models.Gallery, error) {
	return s.GetGallery_(input)
}

func (s *MockGalleryStore) GetGalleries(input *GetGalleriesInput) (*PaginatedGalleriesResult, error) {
	return s.GetGalleries_(input)
}

func (s *MockGalleryStore) GalleryExists(input *GetGalleryInput) (bool, error) {
	return s.GalleryExists_(input)
}

type MockResizer struct {
	ResizeImage_ func(newSize *Size, directory, filename string) ([]byte, error)
}

func (m *MockResizer) ResizeImage(newSize *Size, directory, filename string) ([]byte, error) {
	return m.ResizeImage_(newSize, directory, filename)
}
