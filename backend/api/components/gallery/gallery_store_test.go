package gallery

import (
	"api"
	"db"
	"db/models"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
)

func TestGetGallery(t *testing.T) {
	t.Run("Returns gallery with files by id", func(t *testing.T) {
		tx := db.GetTestTransaction()
		defer tx.Rollback()
		populator := db.NewPopulator(tx)
		store := NewGalleryStore(tx, api.EnvMediaDir.Load())

		files := []*models.File{
			{Filename: "a"},
			{Filename: "b"},
		}
		gallery := populator.Gallery(map[string]any{"Files": files})
		populator.Gallery()

		savedGallery, err := store.GetGallery(&GetGalleryInput{Id: &gallery.Id})

		assert.Nil(t, err)
		assert.Equal(t, len(gallery.Files), len(savedGallery.Files))
		assert.Equal(t, gallery.Title, savedGallery.Title)
	})

	t.Run("Returns gallery with files by slug", func(t *testing.T) {
		tx := db.GetTestTransaction()
		defer tx.Rollback()
		populator := db.NewPopulator(tx)
		store := NewGalleryStore(tx, api.EnvMediaDir.Load())

		files := []*models.File{
			{Filename: "a"},
			{Filename: "b"},
		}
		gallery := populator.Gallery(map[string]any{"Files": files})
		populator.Gallery()

		savedGallery, err := store.GetGallery(&GetGalleryInput{Slug: &gallery.Slug})
		assert.Nil(t, err)
		assert.Equal(t, gallery.Id, savedGallery.Id)
		assert.Equal(t, len(gallery.Files), len(savedGallery.Files))
	})

	t.Run("Returns error if nor Id nor Slug was provided", func(t *testing.T) {
		tx := db.GetTestTransaction()
		defer tx.Rollback()
		store := NewGalleryStore(tx, api.EnvMediaDir.Load())

		_, err := store.GetGallery(&GetGalleryInput{})
		assert.Error(t, err)
	})

	t.Run("Returns NOT_FOUND error if gallery doesn't exist", func(t *testing.T) {
		tx := db.GetTestTransaction()
		defer tx.Rollback()
		populator := db.NewPopulator(tx)
		store := NewGalleryStore(tx, api.EnvMediaDir.Load())

		gallery := populator.Gallery(map[string]any{"Active": false})

		input := &GetGalleryInput{Id: &gallery.Id}
		input.SetActive(true)

		_, err := store.GetGallery(input)

		assert.EqualError(t, err, ErrorNotExisits.Error())
	})
}

func TestGalleryExists(t *testing.T) {
	tx := db.GetTestTransaction()
	defer tx.Rollback()
	populator := db.NewPopulator(tx)
	store := NewGalleryStore(tx, api.EnvMediaDir.Load())

	title := "Gallery1"
	exists, err := store.GalleryExists(&GetGalleryInput{Title: &title})
	assert.Nil(t, err)
	assert.False(t, exists)

	populator.Gallery(map[string]any{
		"Title": *&title,
	})

	exists, err = store.GalleryExists(&GetGalleryInput{Title: &title})
	assert.Nil(t, err)
	assert.True(t, exists)
}

func TestCreateGallery(t *testing.T) {
	t.Run("Returns error if something unexpected happens", func(t *testing.T) {
		tx := db.GetTestTransaction()
		defer tx.Rollback()
		store := NewGalleryStore(tx, api.EnvMediaDir.Load())

		_, err := store.CreateGallery(&CreateUpdateGalleryInput{})

		assert.Error(t, err)
		assert.False(t, strings.Contains(err.Error(), ErrorNotExisits.Error()))
	})

	t.Run("Creates gallery with proper data", func(t *testing.T) {
		tx := db.GetTestTransaction()
		defer tx.Rollback()
		store := NewGalleryStore(tx, api.EnvMediaDir.Load())

		today := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
		date := datatypes.Date(today)
		input := &CreateUpdateGalleryInput{
			Title:       "Title",
			Description: "desc",
			Date:        &date,
			Active:      true,
			Slug:        "slug",
		}

		result, err := store.CreateGallery(input)

		savedGallery := &models.Gallery{}
		tx.Where("title = ?", input.Title).Find(savedGallery)

		assert.Nil(t, err)
		assert.Equal(t, input.Title, result.Title)
		assert.Equal(t, input.Description, result.Description)
		assert.Equal(t, input.Date, result.Date)
		assert.Equal(t, input.Active, result.Active)
		assert.Equal(t, input.Slug, result.Slug)
		assert.NotEqual(t, "", result.Directory)

		assert.Equal(t, input.Title, savedGallery.Title)
		assert.Equal(t, input.Description, savedGallery.Description)
		assert.Equal(t, input.Date, savedGallery.Date)
		assert.Equal(t, input.Active, savedGallery.Active)
		assert.Equal(t, input.Slug, savedGallery.Slug)
		assert.NotEqual(t, "", savedGallery.Directory)
	})
}

func TestUpdateGallery(t *testing.T) {
	t.Run("Updates gallery with proper data", func(t *testing.T) {
		tx := db.GetTestTransaction()
		defer tx.Rollback()
		populator := db.NewPopulator(tx)
		store := NewGalleryStore(tx, api.EnvMediaDir.Load())

		gallery := populator.Gallery()
		today := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
		date := datatypes.Date(today)
		input := &CreateUpdateGalleryInput{
			Title:       "asd",
			Description: "sfdsfg",
			Date:        &date,
			Active:      true,
			Slug:        "slug",
		}

		err := store.UpdateGallery(gallery.Id, input)

		savedGallery := &models.Gallery{}
		tx.Where("title = ?", input.Title).Find(savedGallery)

		assert.Nil(t, err)
		assert.Equal(t, input.Title, savedGallery.Title)
		assert.Equal(t, input.Description, savedGallery.Description)
		assert.Equal(t, input.Date, savedGallery.Date)
		assert.Equal(t, input.Active, savedGallery.Active)
		assert.Equal(t, input.Slug, savedGallery.Slug)
	})
}

func TestDeleteGallery(t *testing.T) {
	t.Run("Returns error if gallery doesn't exist", func(t *testing.T) {
		tx := db.GetTestTransaction()
		defer tx.Rollback()

		store := NewGalleryStore(tx, api.EnvMediaDir.Load())

		err := store.DeleteGallery(&models.Gallery{})

		assert.EqualError(t, ErrorNotExisits, err.Error())
	})

	t.Run("Deletes gallery record", func(t *testing.T) {
		tx := db.GetTestTransaction()
		defer tx.Rollback()
		populator := db.NewPopulator(tx)
		store := NewGalleryStore(tx, api.EnvMediaDir.Load())

		gallery := populator.Gallery()

		err := store.DeleteGallery(gallery)
		assert.Nil(t, err)

		result := int64(0)
		tx.Find(gallery).Count(&result)
		assert.Equal(t, int64(0), result)
	})

	t.Run("Removes all files", func(t *testing.T) {
		// given
		tx := db.GetTestTransaction()
		defer tx.Rollback()
		populator := db.NewPopulator(tx)
		store := NewGalleryStore(tx, api.EnvMediaDir.Load())

		files := []*models.File{
			{Filename: "file1.jpg"},
			{Filename: "file2.jpg"},
		}

		gallery := populator.Gallery(map[string]any{"Files": files})

		basePath := filepath.Join(api.EnvMediaDir.Load(), gallery.Directory)
		if fileErr := os.MkdirAll(basePath, 0755); fileErr != nil {
			panic(fileErr)
		}
		defer os.RemoveAll(basePath)

		file, fileErr := os.Create(filepath.Join(basePath, files[0].Filename))
		if fileErr != nil {
			panic(fileErr)
		}
		file.Close()

		if gallery.Id == 0 || files[0].Id == 0 || files[1].Id == 0 {
			panic("Records weren't created")
		}

		// when
		err := store.DeleteGallery(gallery)

		// then
		assert.Nil(t, err)

		result := int64(0)
		tx.Find(gallery).Count(&result)
		assert.Equal(t, int64(0), result)

		tx.Find(&models.File{}).Count(&result)
		assert.Equal(t, int64(0), result)

		_, notExistsError := os.Stat(basePath)
		assert.True(t, os.IsNotExist(notExistsError))
	})
}

func TestGetGalleries(t *testing.T) {

	t.Run("Returns galleries ordered by date", func(t *testing.T) {
		tx := db.GetTestTransaction()
		defer tx.Rollback()
		populator := db.NewPopulator(tx)
		store := NewGalleryStore(tx, api.EnvMediaDir.Load())

		populator.Gallery()
		populator.Gallery()

		savedGalleries, err := store.GetGalleries(&GetGalleriesInput{})

		assert.Nil(t, err)
		assert.Equal(t, 2, len(savedGalleries.Items))
	})

	t.Run("Returns properly ordered files", func(t *testing.T) {
		tx := db.GetTestTransaction()
		defer tx.Rollback()
		populator := db.NewPopulator(tx)
		store := NewGalleryStore(tx, api.EnvMediaDir.Load())

		files := []*models.File{
			{Filename: "moon.jpg", Rank: 1},
			{Filename: "banana.jpg", Rank: 0},
			{Filename: "sun.jpg", Rank: 2},
			{Filename: "apple.jpg", Rank: 0},
		}

		populator.Gallery(map[string]any{"Files": files})

		savedGalleries, err := store.GetGalleries(&GetGalleriesInput{})
		assert.Nil(t, err)
		assert.Equal(t, len(files), len(savedGalleries.Items[0].Files))
		savedFiles := savedGalleries.Items[0].Files

		assert.Equal(t, files[2].Filename, savedFiles[0].Filename)
		assert.Equal(t, files[0].Filename, savedFiles[1].Filename)
		assert.Equal(t, files[3].Filename, savedFiles[2].Filename)
	})

	t.Run("Filters active", func(t *testing.T) {
		tx := db.GetTestTransaction()
		defer tx.Rollback()
		populator := db.NewPopulator(tx)
		store := NewGalleryStore(tx, api.EnvMediaDir.Load())

		activeGallery := populator.Gallery(map[string]any{"Active": true})
		populator.Gallery(map[string]any{"Active": false})

		input := &GetGalleriesInput{}
		input.SetActive(true)

		savedGalleries, err := store.GetGalleries(input)

		assert.Nil(t, err)
		assert.Equal(t, 1, len(savedGalleries.Items))
		assert.Equal(t, activeGallery.Title, savedGalleries.Items[0].Title)
	})
}
