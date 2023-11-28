package gallery

import (
	"api"
	"db"
	"db/models"
	"testing"

	"github.com/stretchr/testify/assert"
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

// func TestCreateGallery(t *testing.T) {
// 	t.Run("Returns error if something unexpected happens", func(t *testing.T) {
// 		tx := db.GetTestTransaction()
// 		defer tx.Rollback()
// 		store := NewGalleryStore(db, api.EnvMediaDir.Load())

// 		_, err := store.CreateGallery(&CreateUpdateGalleryInput{})

// 		assert.Error(t, err)
// 		assert.False(t, strings.Contains(err.Error(), components.NotFoundError.Error()))
// 	})

// 	t.Run("Creates gallery with proper data", func(t *testing.T) {
// 		tx := db.GetTestTransaction()
// 		defer tx.Rollback()
// 		store := NewGalleryStore(db, api.EnvMediaDir.Load())

// 		wantedGallery := NewDummyGallery()

// 		result, err := store.CreateGallery(&CreateUpdateGalleryInput{
// 			Title:       wantedGallery.Title,
// 			Description: wantedGallery.Description,
// 			Date:        wantedGallery.Date,
// 			Active:      wantedGallery.Active,
// 			Slug:        wantedGallery.Slug,
// 		})

// 		savedGallery := &models.Gallery{}
// 		db.Where("title = ?", wantedGallery.Title).Find(savedGallery)

// 		assert.Nil(t, err)
// 		assert.Equal(t, wantedGallery.Title, result.Title)
// 		assert.Equal(t, wantedGallery.Description, result.Description)
// 		assert.Equal(t, wantedGallery.Date, result.Date)
// 		assert.Equal(t, wantedGallery.Active, result.Active)
// 		assert.Equal(t, wantedGallery.Slug, result.Slug)
// 		assert.NotEqual(t, "", result.Directory)

// 		assert.Equal(t, wantedGallery.Title, savedGallery.Title)
// 		assert.Equal(t, wantedGallery.Description, savedGallery.Description)
// 		assert.Equal(t, wantedGallery.Date, savedGallery.Date)
// 		assert.Equal(t, wantedGallery.Active, savedGallery.Active)
// 		assert.Equal(t, wantedGallery.Slug, savedGallery.Slug)
// 		assert.NotEqual(t, "", savedGallery.Directory)
// 	})
// }

// func TestUpdateGallery(t *testing.T) {
// 	t.Run("Updates gallery with proper data", func(t *testing.T) {
// 		tx := db.GetTestTransaction()
// 		defer tx.Rollback()
// 		store := NewGalleryStore(db, api.EnvMediaDir.Load())

// 		wantedGallery := NewDummyGallery()
// 		initialGallery := NewDummyGallery()
// 		initialGallery.Directory = wantedGallery.Directory

// 		panicErrorResult(db.Create(initialGallery))
// 		wantedGallery.Id = initialGallery.Id
// 		panicErrorResult(db.Create(NewDummyGallery()))

// 		err := store.UpdateGallery(initialGallery.Id, &CreateUpdateGalleryInput{
// 			Title:       wantedGallery.Title,
// 			Description: wantedGallery.Description,
// 			Date:        wantedGallery.Date,
// 			Active:      wantedGallery.Active,
// 			Slug:        wantedGallery.Slug,
// 		})

// 		savedGallery := &models.Gallery{}
// 		db.Where("title = ?", wantedGallery.Title).Find(savedGallery)

// 		assert.Nil(t, err)
// 		assert.Equal(t, wantedGallery, savedGallery)
// 	})

// 	t.Run("Returns error when table doesn't exist", func(t *testing.T) {
// 		db := getEmptyDb()
// 		store := NewGalleryStore(db, api.EnvMediaDir.Load())

// 		err := store.UpdateGallery(1, &CreateUpdateGalleryInput{})
// 		assert.Error(t, err)
// 	})
// }

// func TestDeleteGallery(t *testing.T) {
// 	t.Run("Returns error if gallery doesn't exist", func(t *testing.T) {
// 		tx := db.GetTestTransaction()
// 		defer tx.Rollback()

// 		store := NewGalleryStore(db, api.EnvMediaDir.Load())

// 		err := store.DeleteGallery(&models.Gallery{})

// 		assert.EqualError(t, components.NotFoundError, err.Error())
// 	})

// 	t.Run("Returns error if something unexpected happens", func(t *testing.T) {
// 		db := getEmptyDb()
// 		store := NewGalleryStore(db, api.EnvMediaDir.Load())

// 		err := store.DeleteGallery(&models.Gallery{})

// 		assert.Error(t, err)
// 		assert.False(t, strings.Contains(err.Error(), components.NotFoundError.Error()))
// 	})

// 	t.Run("Deletes gallery record", func(t *testing.T) {
// 		tx := db.GetTestTransaction()
// 		defer tx.Rollback()
// 		store := NewGalleryStore(db, api.EnvMediaDir.Load())

// 		gallery := NewDummyGallery()
// 		db.Create(gallery)

// 		err := store.DeleteGallery(gallery)
// 		assert.Nil(t, err)
// 		result := int64(0)
// 		db.Find(gallery).Count(&result)
// 		assert.Equal(t, int64(0), result)
// 	})

// 	t.Run("Removes all files", func(t *testing.T) {
// 		// given
// 		tx := db.GetTestTransaction()
// 		defer tx.Rollback()
// 		store := NewGalleryStore(db, api.EnvMediaDir.Load())

// 		files := []*models.File{
// 			{Filename: "file1.jpg"},
// 			{Filename: "file2.jpg"},
// 		}

// 		gallery := NewDummyGallery()
// 		gallery.Files = files

// 		basePath := filepath.Join(api.EnvMediaDir.Load(), gallery.Directory)
// 		if fileErr := os.MkdirAll(basePath, 0755); fileErr != nil {
// 			panic(fileErr)
// 		}
// 		defer os.RemoveAll(basePath)

// 		file, fileErr := os.Create(filepath.Join(basePath, files[0].Filename))
// 		if fileErr != nil {
// 			panic(fileErr)
// 		}
// 		file.Close()

// 		db.Create(gallery)

// 		if gallery.Id == 0 || files[0].Id == 0 || files[1].Id == 0 {
// 			panic("Records weren't created")
// 		}

// 		// when
// 		err := store.DeleteGallery(gallery)

// 		// then
// 		assert.Nil(t, err)

// 		result := int64(0)
// 		db.Find(gallery).Count(&result)
// 		assert.Equal(t, int64(0), result)

// 		db.Find(&models.File{}).Count(&result)
// 		assert.Equal(t, int64(0), result)

// 		_, notExistsError := os.Stat(basePath)
// 		assert.True(t, os.IsNotExist(notExistsError))
// 	})
// }

// func TestGetGalleries(t *testing.T) {
// 	t.Run("Returns error if something unexpected happens", func(t *testing.T) {
// 		db := getEmptyDb()
// 		store := NewGalleryStore(db, api.EnvMediaDir.Load())

// 		_, err := store.GetGalleries(&GetGalleriesInput{})

// 		assert.Error(t, err)
// 		assert.False(t, strings.Contains(err.Error(), components.NotFoundError.Error()))
// 	})

// 	t.Run("Returns galleries ordered by date", func(t *testing.T) {
// 		tx := db.GetTestTransaction()
// 		defer tx.Rollback()
// 		store := NewGalleryStore(db, api.EnvMediaDir.Load())

// 		galleries := []*models.Gallery{NewDummyGallery(), NewDummyGallery()}
// 		galleries[0].Date = helpers.DateFromISO8601Panic("2022-01-01")
// 		galleries[1].Date = helpers.DateFromISO8601Panic("2022-01-02")

// 		panicErrorResult(db.Create(&galleries))

// 		savedGalleries, err := store.GetGalleries(&GetGalleriesInput{})

// 		assert.Nil(t, err)
// 		assert.Equal(t, galleries[1].Title, savedGalleries.Items[0].Title)
// 	})

// 	t.Run("Returns properly ordered files", func(t *testing.T) {
// 		tx := db.GetTestTransaction()
// 		defer tx.Rollback()
// 		store := NewGalleryStore(db, api.EnvMediaDir.Load())

// 		galleries := []*models.Gallery{NewDummyGallery()}

// 		galleries[0].Files = []*models.File{
// 			{Filename: "moon.jpg", Rank: 1},
// 			{Filename: "banana.jpg", Rank: 0},
// 			{Filename: "sun.jpg", Rank: 2},
// 			{Filename: "apple.jpg", Rank: 0},
// 		}

// 		panicErrorResult(db.Create(&galleries))

// 		savedGalleries, err := store.GetGalleries(&GetGalleriesInput{})
// 		assert.Nil(t, err)
// 		assert.Equal(t, len(galleries[0].Files), len(savedGalleries.Items[0].Files))
// 		files := galleries[0].Files
// 		savedFiles := savedGalleries.Items[0].Files

// 		assert.Equal(t, files[2].Filename, savedFiles[0].Filename)
// 		assert.Equal(t, files[0].Filename, savedFiles[1].Filename)
// 		assert.Equal(t, files[3].Filename, savedFiles[2].Filename)
// 	})

// 	t.Run("Filters active", func(t *testing.T) {
// 		tx := db.GetTestTransaction()
// 		defer tx.Rollback()
// 		store := NewGalleryStore(db, api.EnvMediaDir.Load())

// 		galleries := []*models.Gallery{NewDummyGallery(), NewDummyGallery()}
// 		galleries[0].Active = true
// 		galleries[1].Active = false

// 		panicErrorResult(db.Create(&galleries))
// 		input := &GetGalleriesInput{}
// 		input.SetActive(true)

// 		savedGalleries, err := store.GetGalleries(input)

// 		assert.Nil(t, err)
// 		assert.Equal(t, 1, len(savedGalleries.Items))
// 		assert.Equal(t, galleries[0].Title, savedGalleries.Items[0].Title)
// 	})
// }
