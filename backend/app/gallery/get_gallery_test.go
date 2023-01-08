package app_gallery

import (
	"attiladudas/backend/components"
	"attiladudas/backend/components/gallery"
	"attiladudas/backend/helpers"
	"attiladudas/backend/models"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetGallery(t *testing.T) {
	path := func(slug string) string {
		return fmt.Sprintf("/api/galleries/%s/", slug)
	}

	t.Run("Returns gallery", func(t *testing.T) {
		app := newMockApp()

		files := []*models.File{
			{Filename: "1.jpg", Rank: 3},
		}
		galleryObj := gallery.NewDummyGallery()
		galleryObj.Files = files

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
			Date:        helpers.DateToISO8601(*galleryObj.Date),
			Description: galleryObj.Description,
			Active:      galleryObj.Active,
			Files:       filesResponse,
		}

		app.galleryStore.GetGallery_ = func(input *gallery.GetGalleryInput) (*models.Gallery, error) {
			if *input.Slug == galleryObj.Slug && *input.Active {
				return galleryObj, nil
			}
			return nil, components.NotFoundError
		}

		app.fileStore.GetFileDownloadPath_ = func(gallery *models.Gallery, filename string) string {
			return wantedResponse.Files[0].Path
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(
			"GET",
			path(galleryObj.Slug),
			&bytes.Buffer{},
		)
		query := req.URL.Query()
		req.URL.RawQuery = query.Encode()

		app.engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		response := &GalleryResponse{}
		err := json.Unmarshal(w.Body.Bytes(), response)
		assert.Nil(t, err)
		assert.Equal(t, wantedResponse, response)
	})

	t.Run("Returns INTERNAL SERVER ERROR on store error", func(t *testing.T) {
		app := newMockApp()
		app.galleryStore.GetGallery_ = func(input *gallery.GetGalleryInput) (*models.Gallery, error) {
			return nil, errors.New("ARMAGEDDON")
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(
			"GET",
			path("asd"),
			&bytes.Buffer{},
		)

		app.engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Returns NOT_FOUND when gallery does not exist", func(t *testing.T) {
		app := newMockApp()
		app.galleryStore.GetGallery_ = func(input *gallery.GetGalleryInput) (*models.Gallery, error) {
			return nil, components.NotFoundError
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(
			"GET",
			path("asd"),
			&bytes.Buffer{},
		)

		app.engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
