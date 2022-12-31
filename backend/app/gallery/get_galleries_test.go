package app_gallery

import (
	"attiladudas/backend/components/auth"
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

func TestGetGalleries(t *testing.T) {
	path := "/api/galleries/"

	t.Run("Returns galleries", func(t *testing.T) {
		app := newMockApp()

		dateString := "2022-01-01"
		galleries := []*models.Gallery{
			{
				Id:          1,
				Title:       "Gallery1",
				Description: "desc1",
				Date:        helpers.DateFromISO8601Panic(dateString),
				Active:      true,
				Files: []*models.File{
					{Filename: "1.jpg", Rank: 3},
				},
			},
			{
				Id:    2,
				Title: "Gallery2",
				Date:  helpers.DateFromISO8601Panic(dateString),
			},
		}

		wantedResponse := &GalleriesResponse{
			Galleries: []*GalleryResponse{
				{
					Id:          galleries[0].Id,
					Title:       galleries[0].Title,
					Description: galleries[0].Description,
					Active:      galleries[0].Active,
					Date:        dateString,
					Files: []*File{
						{
							Filename: galleries[0].Files[0].Filename,
							Path:     "fake path",
							Rank:     galleries[0].Files[0].Rank,
						},
					},
				},
				{
					Id:    galleries[1].Id,
					Title: galleries[1].Title,
					Files: []*File{},
					Date:  dateString,
				},
			},
			Total: 2,
		}

		wantedInput := &gallery.GetGalleriesInput{
			Page:     2,
			PageSize: 5,
		}

		app.galleryStore.GetGalleries_ = func(input *gallery.GetGalleriesInput) (*gallery.PaginatedGalleriesResult, error) {
			assert.Equal(t, wantedInput, input)
			return &gallery.PaginatedGalleriesResult{Items: galleries, Total: uint(len(wantedResponse.Galleries))}, nil
		}

		app.fileStore.GetFileDownloadPath_ = func(gallery *models.Gallery, filename string) string {
			return wantedResponse.Galleries[0].Files[0].Path
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(
			"GET",
			path,
			&bytes.Buffer{},
		)
		query := req.URL.Query()
		query.Add("page_size", fmt.Sprint(wantedInput.PageSize))
		query.Add("page", fmt.Sprint(wantedInput.Page))
		req.URL.RawQuery = query.Encode()

		app.engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		response := &GalleriesResponse{}
		err := json.Unmarshal(w.Body.Bytes(), response)
		assert.Nil(t, err)
		assert.Equal(t, wantedResponse, response)
	})

	t.Run("Returns INTERNAL SERVER ERROR on store error", func(t *testing.T) {
		app := newMockApp()
		app.galleryStore.GetGalleries_ = func(input *gallery.GetGalleriesInput) (*gallery.PaginatedGalleriesResult, error) {
			return nil, errors.New("ARMAGEDDON")
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(
			"GET",
			path,
			&bytes.Buffer{},
		)
		query := req.URL.Query()
		query.Add("page_size", "2")
		query.Add("page", "2")
		req.URL.RawQuery = query.Encode()

		app.engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Returns BAD REQUEST on invalid query params", func(t *testing.T) {
		app := newMockApp()
		app.galleryStore.GetGalleries_ = func(input *gallery.GetGalleriesInput) (*gallery.PaginatedGalleriesResult, error) {
			return nil, errors.New("ARMAGEDDON")
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(
			"GET",
			path,
			&bytes.Buffer{},
		)
		query := req.URL.Query()
		query.Add("page_size", "-2")
		req.URL.RawQuery = query.Encode()

		app.engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Returns only active galeries when user is not logged in", func(t *testing.T) {
		app := newMockApp()
		app.authContext.RequireUsername_ = func(authHeader string, i auth.IJwt) error {
			return errors.New("missing auth header")
		}

		app.galleryStore.GetGalleries_ = func(input *gallery.GetGalleriesInput) (*gallery.PaginatedGalleriesResult, error) {
			assert.True(t, *input.Active)
			return &gallery.PaginatedGalleriesResult{}, nil
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(
			"GET",
			path,
			&bytes.Buffer{},
		)
		query := req.URL.Query()
		req.URL.RawQuery = query.Encode()

		app.engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
