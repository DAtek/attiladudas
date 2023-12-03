package gallery

import (
	"api"
	"db"
	"db/models"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
 * A running postgres instance is required to run these tests
 */
func TestAddFiles(t *testing.T) {
	mediaDir := api.EnvMediaDir.Load()

	t.Run("AddFiles saves files properly", func(t *testing.T) {
		tx := db.GetTestTransaction()
		defer tx.Rollback()
		populator := db.NewPopulator(tx)
		store := NewFileStore(tx, api.EnvMediaDir.Load())
		defer os.RemoveAll(mediaDir)

		files := []*FileInput{
			{
				Content:  []byte("apple"),
				Filename: "apple.jpg",
			},
			{
				Content:  []byte("lemon"),
				Filename: "lemon.jpg",
			},
		}

		gallery := populator.Gallery()

		err := store.AddFiles(gallery, files)
		assert.Nil(t, err)

		count := int64(0)
		tx.Model(&models.File{}).Where("filename = ?", files[0].Filename).Count(&count)
		assert.Equal(t, int64(1), count)

		basePath := filepath.Join(mediaDir, gallery.Directory)
		for _, file := range files {
			path := filepath.Join(basePath, file.Filename)
			content, err := os.ReadFile(path)
			assert.Nil(t, err)
			assert.Equal(t, file.Content, content)
		}
	})

	t.Run("AddFiles returns error if file already exists", func(t *testing.T) {
		tx := db.GetTestTransaction()
		defer tx.Rollback()
		populator := db.NewPopulator(tx)
		store := NewFileStore(tx, api.EnvMediaDir.Load())
		defer os.RemoveAll(mediaDir)

		gallery := populator.Gallery(map[string]any{
			"Title":     "gallery1",
			"Directory": "gallery1",
			"Files": []*models.File{
				{Filename: "oi.mate"},
			},
		})

		files := []*FileInput{
			{Filename: gallery.Files[0].Filename},
		}

		err := store.AddFiles(gallery, files)

		assert.Error(t, err)
	})
}

func TestUpdateFileRank(t *testing.T) {
	t.Run("Updates rank properly", func(t *testing.T) {
		tx := db.GetTestTransaction()
		defer tx.Rollback()
		populator := db.NewPopulator(tx)
		store := NewFileStore(tx, api.EnvMediaDir.Load())

		files := []*models.File{
			{Filename: "file1.jpg", Rank: 2},
		}

		populator.Gallery(map[string]any{
			"Files": files,
		})

		wantedFile := &models.File{Id: files[0].Id, Rank: -3, Filename: files[0].Filename, GalleryId: files[0].GalleryId}
		err := store.UpdateFileRank(&UpdateFileRankInput{FileId: wantedFile.Id, Rank: wantedFile.Rank})

		assert.Nil(t, err)

		savedFile := &models.File{}
		tx.Where("id = ?", files[0].Id).Find(savedFile)

		assert.Equal(t, wantedFile, savedFile)
	})

	t.Run("Returns error when table doesn't exist", func(t *testing.T) {
		tx := db.GetTestTransaction()
		defer tx.Rollback()
		store := NewFileStore(tx, api.EnvMediaDir.Load())
		err := store.UpdateFileRank(&UpdateFileRankInput{})
		assert.Error(t, err)
	})
}

func TestGetFile(t *testing.T) {
	t.Run("Returns file", func(t *testing.T) {
		tx := db.GetTestTransaction()
		defer tx.Rollback()
		populator := db.NewPopulator(tx)
		store := NewFileStore(tx, api.EnvMediaDir.Load())

		files := []*models.File{
			{Filename: "file1.jpg", Rank: 2},
			{Filename: "file2.jpg", Rank: 3},
		}

		populator.Gallery(map[string]any{"Files": files})

		savedFile, err := store.GetFile(files[0].Id)

		assert.Nil(t, err)
		assert.Equal(t, files[0], savedFile)
	})
}

func TestDeleteFiles(t *testing.T) {
	t.Run("Deletes all files", func(t *testing.T) {
		tx := db.GetTestTransaction()
		defer tx.Rollback()
		populator := db.NewPopulator(tx)
		mediaDir := api.EnvMediaDir.Load()
		store := NewFileStore(tx, mediaDir)
		defer os.RemoveAll(mediaDir)

		// given
		files := []*models.File{
			{Filename: "file1.jpg"},
			{Filename: "file2.jpg"},
			{Filename: "file3.jpg"},
		}

		gallery := populator.Gallery(map[string]any{"Files": files})

		basePath := filepath.Join(mediaDir, gallery.Directory)
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
		idsToDelete := []uint{files[0].Id, files[1].Id}
		err := store.DeleteFiles(gallery, idsToDelete)

		// then
		assert.Nil(t, err)
		count := int64(0)
		tx.Model(&models.File{}).Where("id IN ?", idsToDelete).Count(&count)
		assert.Equal(t, int64(0), count)
		for _, file := range []*models.File{files[0], files[1]} {
			_, err := os.Stat(filepath.Join(basePath, file.Filename))
			assert.True(t, os.IsNotExist(err))
		}
	})
}

func TestMediaDirName(t *testing.T) {
	tx := db.GetTestTransaction()
	defer tx.Rollback()
	store := NewFileStore(tx, "/home/app/files")

	assert.Equal(t, "files", store.MediaDirName())
}
