package app_gallery

import (
	"attiladudas/backend/components/auth"
	"attiladudas/backend/components/gallery"
	"attiladudas/backend/helpers"
	"attiladudas/backend/models"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostGallery(t *testing.T) {
	path := "/api/gallery/"
	validInput := createUpdateGalleryBody{
		Title:       "Gallery1",
		Slug:        "gallery-1",
		Description: "Desc1",
		Date:        "2020-01-01",
		Active:      true,
	}

	t.Run("Test unauthorized request", func(t *testing.T) {
		app := newMockApp()

		app.authContext.RequireUsername_ = func(authHeader string, i auth.IJwt) error {
			return auth.InvalidAuthHeaderError
		}

		w := httptest.NewRecorder()
		body, _ := json.Marshal(gallery.CreateUpdateGalleryInput{})
		req, _ := http.NewRequest(
			"POST",
			path,
			bytes.NewBuffer([]byte(body)),
		)
		app.engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	badDate, _ := json.Marshal(createUpdateGalleryBody{
		Title:       "Gallery1",
		Description: "Desc1",
		Date:        "asd",
		Active:      true,
	})
	badRequestBodies := [][]byte{
		[]byte(""),
		badDate,
	}
	for _, body := range badRequestBodies {
		t.Run("Bad request if body is invalid", func(t *testing.T) {
			app := newMockApp()

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(
				"POST",
				path,
				bytes.NewBuffer(body),
			)
			app.engine.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)
		})
	}

	t.Run("Bad request ifgallery with slug exists", func(t *testing.T) {
		app := newMockApp()

		app.galleryStore.GalleryExists_ = func(input *gallery.GetGalleryInput) (bool, error) {
			return true, nil
		}

		body, _ := json.Marshal(validInput)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(
			"POST",
			path,
			bytes.NewBuffer(body),
		)
		app.engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Internal server error if creation fails", func(t *testing.T) {
		app := newMockApp()
		app.jwtContext.Decode_ = func(s string) (*auth.Claims, error) {
			return &auth.Claims{Username: "user1"}, nil
		}
		w := httptest.NewRecorder()
		body, _ := json.Marshal(validInput)
		app.galleryStore.CreateGallery_ = func(cugi *gallery.CreateUpdateGalleryInput) (*models.Gallery, error) {
			return nil, errors.New("Already exists")
		}
		req, _ := http.NewRequest(
			"POST",
			path,
			bytes.NewBuffer(body),
		)
		app.engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Creates gallery", func(t *testing.T) {
		app := newMockApp()
		app.jwtContext.Decode_ = func(s string) (*auth.Claims, error) {
			return &auth.Claims{Username: "user1"}, nil
		}
		createCalled := false

		// Calls create gallery with proper arguments
		app.galleryStore.CreateGallery_ = func(i *gallery.CreateUpdateGalleryInput) (*models.Gallery, error) {
			createCalled = true
			wantedInput := &gallery.CreateUpdateGalleryInput{
				Title:       validInput.Title,
				Slug:        validInput.Slug,
				Description: validInput.Description,
				Date:        helpers.DateFromISO8601Panic(validInput.Date),
				Active:      validInput.Active,
			}
			assert.Equal(t, wantedInput, i)
			return &models.Gallery{Id: 1}, nil
		}

		w := httptest.NewRecorder()

		body, _ := json.Marshal(validInput)

		req, _ := http.NewRequest(
			"POST",
			path,
			bytes.NewBuffer(body),
		)
		app.engine.ServeHTTP(w, req)
		response := createGalleryResponse{}
		decoder := json.NewDecoder(w.Body)
		err := decoder.Decode(&response)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Nil(t, err)
		assert.Equal(t, uint(1), response.Id)
		assert.True(t, createCalled)
	})

	t.Run("Bad request if gallery with same title already exists", func(t *testing.T) {
		app := newMockApp()

		app.galleryStore.GalleryExists_ = func(input *gallery.GetGalleryInput) (bool, error) {
			if input.Title != nil {
				return *input.Title == validInput.Title, nil
			}

			return false, nil
		}

		w := httptest.NewRecorder()
		body, _ := json.Marshal(validInput)
		req, _ := http.NewRequest(
			"POST",
			path,
			bytes.NewBuffer(body),
		)
		app.engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		err := &helpers.JsonErrorCollection{}
		json.Unmarshal(w.Body.Bytes(), err)

		assert.Equal(t, 1, len(err.Errors))
		assert.Equal(t, helpers.ErrorAlreadyExists, *err.Errors[0].Type)
		assert.Equal(t, "title", *err.Errors[0].Location)
	})
}
