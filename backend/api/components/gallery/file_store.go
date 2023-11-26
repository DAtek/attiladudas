package gallery

import (
	"attiladudas/backend/models"
	"os"
	"path/filepath"
	"strings"

	"gorm.io/gorm"
)

type IFileStore interface {
	AddFiles(*models.Gallery, []*FileInput) error
	DeleteFiles(*models.Gallery, []uint) error
	GetFileDownloadPath(gallery *models.Gallery, filename string) string
	UpdateFileRank(file *UpdateFileRankInput) error
	GetFile(id uint) (*models.File, error)
	MediaDirName() string
}

type UpdateFileRankInput struct {
	FileId uint
	Rank   int
}

type fileStore struct {
	db      *gorm.DB
	baseDir string
}

func NewFileStore(db *gorm.DB, baseDir string) IFileStore {
	return &fileStore{db: db, baseDir: baseDir}
}

func (store *fileStore) AddFiles(gallery *models.Gallery, files []*FileInput) error {
	dbFiles := []*models.File{}
	for _, file := range files {
		dbFiles = append(dbFiles, &models.File{GalleryId: gallery.Id, Filename: file.Filename})
	}

	if result := store.db.Create(dbFiles); result.Error != nil {
		return result.Error
	}

	basePath := filepath.Join(store.baseDir, gallery.Directory)
	if _, notExistsError := os.Stat(basePath); os.IsNotExist(notExistsError) {
		os.MkdirAll(basePath, os.FileMode(0755))
	}

	for _, file := range files {
		path := store.getFileFullPath(gallery, file.Filename)
		f, fileErr := os.Create(path)
		if fileErr != nil {
			// In theory this should never happen
			return fileErr
		}
		defer f.Close()
		_, fileErr = f.Write(file.Content)
		if fileErr != nil {
			// In theory this should never happen
			return fileErr
		}
	}

	return nil
}

func (store *fileStore) DeleteFiles(gallery *models.Gallery, fileIds []uint) error {
	files := []*models.File{}
	store.db.Where("id IN ?", fileIds).Find(&files)
	for _, file := range files {
		path := store.getFileFullPath(gallery, file.Filename)
		os.Remove(path)
	}
	return store.db.Delete(&models.File{}, "id IN ?", fileIds).Error
}

func (store *fileStore) GetFileDownloadPath(gallery *models.Gallery, filename string) string {
	return filepath.Join("/", store.MediaDirName(), gallery.Directory, filename)
}

func (store *fileStore) MediaDirName() string {
	mediaDirParts := strings.Split(store.baseDir, "/")
	return mediaDirParts[len(mediaDirParts)-1]
}

func (store *fileStore) getFileFullPath(gallery *models.Gallery, filename string) string {
	return filepath.Join(store.baseDir, gallery.Directory, filename)
}

func (store *fileStore) UpdateFileRank(data *UpdateFileRankInput) error {
	result := store.db.Model(&models.File{Id: data.FileId}).Updates(map[string]any{
		"Rank": data.Rank,
	})
	return result.Error
}

func (store *fileStore) GetFile(id uint) (*models.File, error) {
	file := &models.File{}
	result := store.db.Where("id = ?", id).Find(file)
	return file, result.Error
}
