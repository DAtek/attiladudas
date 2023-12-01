package gallery_get

import (
	"api"
	"api/components/gallery"
	"api/helpers"
	"bytes"
	"db/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/DAtek/gotils"
	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
)

func TestGetGallery(t *testing.T) {
	getPath := func(slug string) string {
		return fmt.Sprintf("/api/gallery/%s/", slug)
	}

	t.Run("Returns gallery", func(t *testing.T) {
		files := []*models.File{
			{Filename: "1.jpg", Rank: 3},
		}

		today := time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC)
		date := datatypes.Date(today)
		galleryObj := &models.Gallery{
			Id:          1,
			Title:       "title",
			Slug:        "gallery-1",
			Description: "desc",
			Date:        &date,
			Files:       files,
			Directory:   "dir",
			Active:      true,
		}

		filesResponse := []*File{
			{
				Id:       files[0].Id,
				Filename: files[0].Filename,
				Path:     "fake_path",
				Rank:     files[0].Rank,
			},
		}
		wantedResponse := &GalleryResponse{
			Id:          galleryObj.Id,
			Title:       galleryObj.Title,
			Slug:        galleryObj.Slug,
			Date:        helpers.DateToISO8601(galleryObj.Date),
			Description: galleryObj.Description,
			Active:      galleryObj.Active,
			Files:       filesResponse,
		}

		galleryStore := &gallery.MockGalleryStore{
			GetGallery_: func(input *gallery.GetGalleryInput) (*models.Gallery, error) {
				if *input.Slug == galleryObj.Slug && *input.Active {
					return galleryObj, nil
				}
				return nil, nil
			},
		}

		fileStore := &gallery.MockFileStore{
			GetFileDownloadPath_: func(gallery *models.Gallery, filename string) string {
				return wantedResponse.Files[0].Path
			},
		}

		app := api.AppWithMiddlewares(
			PluginGetGallery(galleryStore, fileStore),
		)

		req := gotils.ResultOrPanic(http.NewRequest(
			"GET",
			getPath(galleryObj.Slug),
			&bytes.Buffer{},
		))

		resp := gotils.ResultOrPanic(app.Test(req))

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		response := &GalleryResponse{}
		respData := &bytes.Buffer{}
		io.Copy(respData, resp.Body)
		json.Unmarshal(respData.Bytes(), response)
		assert.Equal(t, wantedResponse, response)
	})

	t.Run("Returns NOT_FOUND when gallery does not exist", func(t *testing.T) {
		galleryStore := &gallery.MockGalleryStore{
			GetGallery_: func(input *gallery.GetGalleryInput) (*models.Gallery, error) {
				return nil, nil
			},
		}

		fileStore := &gallery.MockFileStore{}

		app := api.AppWithMiddlewares(
			PluginGetGallery(galleryStore, fileStore),
		)

		req := gotils.ResultOrPanic(http.NewRequest(
			"GET",
			getPath("title-1"),
			&bytes.Buffer{},
		))

		resp := gotils.ResultOrPanic(app.Test(req))

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}
