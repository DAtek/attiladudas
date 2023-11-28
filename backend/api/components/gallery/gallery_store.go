package gallery

import (
	"db/models"
	"errors"
	"os"
	"path/filepath"

	"github.com/DAtek/gopaq"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type GalleryError string

func (e GalleryError) Error() string {
	return string(e)
}

const (
	ErrorNotExisits = GalleryError("NOT_EXISTS")
)

type CreateUpdateGalleryInput struct {
	Title       string
	Slug        string
	Description string
	Date        *datatypes.Date
	Active      bool
}

type FileInput struct {
	Content  []byte
	Filename string
}

type GetGalleriesInput struct {
	Page     uint
	PageSize uint
	Active   *bool
}

func (input *GetGalleriesInput) SetActive(value bool) {
	input.Active = &value
}

type GetGalleryInput struct {
	Id     *uint
	Slug   *string
	Title  *string
	Active *bool
}

func (input *GetGalleryInput) SetActive(value bool) {
	input.Active = &value
}

type PaginatedGalleriesResult = gopaq.PaginatedQueryResult[[]*models.Gallery]

type IGalleryStore interface {
	CreateGallery(input *CreateUpdateGalleryInput) (*models.Gallery, error)
	UpdateGallery(galleryId uint, input *CreateUpdateGalleryInput) error
	DeleteGallery(gallery *models.Gallery) error
	GetGallery(input *GetGalleryInput) (*models.Gallery, error)
	GetGalleries(input *GetGalleriesInput) (*PaginatedGalleriesResult, error)
	GalleryExists(input *GetGalleryInput) (bool, error)
}

type galleryStore struct {
	db      *gorm.DB
	baseDir string
}

func NewGalleryStore(db *gorm.DB, baseDir string) IGalleryStore {
	return &galleryStore{db: db, baseDir: baseDir}
}

func (store *galleryStore) GetGallery(input *GetGalleryInput) (*models.Gallery, error) {
	if input.Id == nil && input.Slug == nil {
		return nil, errors.New("INVALID_INPUT")
	}

	gallery := &models.Gallery{}
	query := store.db.Preload("Files")
	query = buildWhereConditions(query, input)
	result := query.Find(gallery)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, ErrorNotExisits
	}

	return gallery, nil
}

func (store *galleryStore) GalleryExists(input *GetGalleryInput) (bool, error) {
	count := int64(0)
	query := store.db.Model(&models.Gallery{})
	query = buildWhereConditions(query, input)
	result := query.Count(&count)
	return count > 0, result.Error
}

func buildWhereConditions(query *gorm.DB, input *GetGalleryInput) *gorm.DB {
	newQuery := query

	if input.Id != nil {
		newQuery = query.Where("id = ?", input.Id)
	}

	if input.Slug != nil {
		newQuery = query.Where("slug = ?", input.Slug)
	}

	if input.Active != nil {
		newQuery = query.Where("active = ?", input.Active)
	}

	if input.Title != nil {
		newQuery = query.Where("title = ?", input.Title)
	}

	return newQuery
}

func (store *galleryStore) CreateGallery(data *CreateUpdateGalleryInput) (*models.Gallery, error) {
	gallery := &models.Gallery{
		Title:       data.Title,
		Description: data.Description,
		Date:        data.Date,
		Active:      data.Active,
		Slug:        data.Slug,
		Directory:   uuid.NewString(),
	}
	result := store.db.Create(gallery)
	return gallery, result.Error
}

func (store *galleryStore) UpdateGallery(galleryId uint, data *CreateUpdateGalleryInput) error {
	result := store.db.Model(&models.Gallery{}).Where("id = ?", galleryId).Updates(map[string]any{
		"Title":       data.Title,
		"Description": data.Description,
		"Date":        data.Date,
		"Active":      data.Active,
		"Slug":        data.Slug,
	})
	return result.Error
}

func (store *galleryStore) DeleteGallery(gallery *models.Gallery) error {
	result := store.db.Where("id = ?", gallery.Id).Delete(gallery)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrorNotExisits
	}

	dir := filepath.Join(store.baseDir, gallery.Directory)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil
	}

	return os.RemoveAll(dir)
}

func (store *galleryStore) GetGalleries(input *GetGalleriesInput) (*PaginatedGalleriesResult, error) {

	query := store.db.Model(&models.Gallery{}).Preload("Files", preloadFiles).Order("date DESC")
	if input.Active != nil {
		query = query.Where("active = ?", input.Active)
	}

	return gopaq.FindWithPagination(
		query,
		[]*models.Gallery{},
		input.Page,
		input.PageSize,
	)
}

func preloadFiles(db *gorm.DB) *gorm.DB {
	return db.Order("file.\nrank\n DESC").Order("file.filename")
}
