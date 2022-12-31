package app_gallery

import (
	"attiladudas/backend/components"
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

func TestPutGallery(t *testing.T) {
	path := func(id any) string {
		return fmt.Sprintf("/api/galleries/%v/", id)
	}

	validInput := createUpdateGalleryBody{
		Title:       "Gallery1",
		Slug:        "gallery-1",
		Description: "Desc1",
		Date:        "2020-01-01",
		Active:      true,
	}
	validId := uint(2)
	validBody, _ := json.Marshal(validInput)

	t.Run("Test unauthorized", func(t *testing.T) {
		app := newMockApp()
		app.authContext.RequireUsername_ = func(authHeader string, i auth.IJwt) error {
			return auth.InvalidAuthHeaderError
		}

		w := httptest.NewRecorder()
		body, _ := json.Marshal(gallery.CreateUpdateGalleryInput{})
		req, _ := http.NewRequest(
			"PUT",
			path(1),
			bytes.NewBuffer([]byte(body)),
		)
		app.engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	badDate, _ := json.Marshal(createUpdateGalleryBody{
		Title:       "Gallery1",
		Description: "Desc1",
		Date:        "1",
		Active:      true,
	})
	badRequestscenarios := [][]byte{
		[]byte(""),
		badDate,
	}
	for _, body := range badRequestscenarios {
		t.Run("Bad request if body is invalid", func(t *testing.T) {
			app := newMockApp()
			app.galleryStore.GetGallery_ = func(input *gallery.GetGalleryInput) (*models.Gallery, error) {
				return &models.Gallery{}, nil
			}
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(
				"PUT",
				path(1),
				bytes.NewBuffer(body),
			)
			app.engine.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)
		})
	}

	t.Run("Bad request if gallery with same title already exists", func(t *testing.T) {
		app := newMockApp()
		app.galleryStore.GetGallery_ = func(input *gallery.GetGalleryInput) (*models.Gallery, error) {
			return &models.Gallery{Title: "asd"}, nil
		}

		app.galleryStore.GalleryExists_ = func(input *gallery.GetGalleryInput) (bool, error) {
			return true, nil
		}

		body, _ := json.Marshal(validInput)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(
			"PUT",
			path(1),
			bytes.NewBuffer(body),
		)
		app.engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Not found if id is invalid", func(t *testing.T) {
		app := newMockApp()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(
			"PUT",
			path("asd"),
			&bytes.Buffer{},
		)
		app.engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Not found if gallery doesn't exist", func(t *testing.T) {
		app := newMockApp()
		app.galleryStore.GetGallery_ = func(input *gallery.GetGalleryInput) (*models.Gallery, error) {
			return nil, components.NotFoundError
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(
			"PUT",
			path(1),
			&bytes.Buffer{},
		)
		app.engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Internal server error if db fails", func(t *testing.T) {
		app := newMockApp()
		w := httptest.NewRecorder()
		app.galleryStore.GetGallery_ = func(input *gallery.GetGalleryInput) (*models.Gallery, error) {
			return nil, errors.New("Oops")
		}
		req, _ := http.NewRequest(
			"PUT",
			path(1),
			&bytes.Buffer{},
		)
		app.engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Internal server error if update fails", func(t *testing.T) {
		app := newMockApp()
		w := httptest.NewRecorder()
		updateCalled := false

		app.galleryStore.GetGallery_ = func(input *gallery.GetGalleryInput) (*models.Gallery, error) {
			return &models.Gallery{}, nil
		}

		app.galleryStore.UpdateGallery_ = func(galleryId uint, i *gallery.CreateUpdateGalleryInput) error {
			updateCalled = true
			return errors.New("Unexpected error")
		}
		req, _ := http.NewRequest(
			"PUT",
			path(1),
			bytes.NewBuffer(validBody),
		)
		app.engine.ServeHTTP(w, req)

		assert.True(t, updateCalled)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Updates gallery", func(t *testing.T) {
		app := newMockApp()
		updateCalled := false

		app.galleryStore.GetGallery_ = func(input *gallery.GetGalleryInput) (*models.Gallery, error) {
			return &models.Gallery{
				Id:    *input.Id,
				Title: validInput.Title,
				Slug:  validInput.Slug,
			}, nil
		}

		app.galleryStore.GalleryExists_ = func(input *gallery.GetGalleryInput) (bool, error) {
			return true, nil
		}

		// Calls update gallery with proper arguments
		app.galleryStore.UpdateGallery_ = func(galleryId uint, i *gallery.CreateUpdateGalleryInput) error {
			updateCalled = true
			assert.Equal(t, validId, galleryId)
			wantedInput := &gallery.CreateUpdateGalleryInput{
				Title:       validInput.Title,
				Slug:        validInput.Slug,
				Description: validInput.Description,
				Date:        helpers.DateFromISO8601Panic(validInput.Date),
				Active:      validInput.Active,
			}
			assert.Equal(t, wantedInput, i)
			return nil
		}

		w := httptest.NewRecorder()

		req, _ := http.NewRequest(
			"PUT",
			path(validId),
			bytes.NewBuffer([]byte(validBody)),
		)
		app.engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.True(t, updateCalled)
	})
}
