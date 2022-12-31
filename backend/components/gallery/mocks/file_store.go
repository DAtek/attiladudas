package gallery_mocks

import (
	"attiladudas/backend/components/gallery"
	"attiladudas/backend/models"
)

type MockFileStore struct {
	AddFiles_            func(*models.Gallery, []*gallery.FileInput) error
	DeleteFiles_         func(*models.Gallery, []uint) error
	GetFileDownloadPath_ func(gallery *models.Gallery, filename string) string
	UpdateFileRank_      func(data *gallery.UpdateFileRankInput) error
	GetFile_             func(id uint) (*models.File, error)
	MediaDirName_        func() string
}

func (s *MockFileStore) AddFiles(gallery *models.Gallery, files []*gallery.FileInput) error {
	return s.AddFiles_(gallery, files)
}

func (s *MockFileStore) DeleteFiles(gallery *models.Gallery, fileIds []uint) error {
	return s.DeleteFiles_(gallery, fileIds)
}

func (s *MockFileStore) GetFileDownloadPath(gallery *models.Gallery, filename string) string {
	return s.GetFileDownloadPath_(gallery, filename)
}

func (s *MockFileStore) UpdateFileRank(data *gallery.UpdateFileRankInput) error {
	return s.UpdateFileRank_(data)
}

func (s *MockFileStore) GetFile(id uint) (*models.File, error) {
	return s.GetFile_(id)
}

func (s *MockFileStore) MediaDirName() string {
	return s.MediaDirName_()
}
