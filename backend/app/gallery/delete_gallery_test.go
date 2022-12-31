package app_gallery

import (
	"attiladudas/backend/components"
	"attiladudas/backend/components/auth"
	"attiladudas/backend/components/gallery"
	"attiladudas/backend/models"
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteGallery(t *testing.T) {
	path := func(id any) string {
		return fmt.Sprintf("/api/galleries/%v/", id)
	}

	t.Run("Test unauthorized", func(t *testing.T) {
		app := newMockApp()

		app.authContext.RequireUsername_ = func(authHeader string, i auth.IJwt) error {
			return auth.InvalidAuthHeaderError
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(
			"DELETE",
			path(1),
			&bytes.Buffer{},
		)
		app.engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Delete gallery OK", func(t *testing.T) {
		app := newMockApp()
		galleryId := uint(2)

		app.galleryStore.GetGallery_ = func(input *gallery.GetGalleryInput) (*models.Gallery, error) {
			return &models.Gallery{Id: galleryId}, nil
		}

		app.galleryStore.DeleteGallery_ = func(gallery *models.Gallery) error {
			assert.Equal(t, galleryId, gallery.Id)
			return nil
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(
			"DELETE",
			path(galleryId),
			&bytes.Buffer{},
		)
		app.engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Not found if gallery not found", func(t *testing.T) {
		app := newMockApp()
		app.galleryStore.GetGallery_ = func(input *gallery.GetGalleryInput) (*models.Gallery, error) {
			return nil, components.NotFoundError
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(
			"DELETE",
			path(1),
			&bytes.Buffer{},
		)
		app.engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Not found if invalid id", func(t *testing.T) {
		app := newMockApp()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(
			"DELETE",
			path("asd"),
			&bytes.Buffer{},
		)
		app.engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Internal server error if something unexpected happens", func(t *testing.T) {
		app := newMockApp()
		app.galleryStore.DeleteGallery_ = func(gallery *models.Gallery) error {
			return errors.New("Earthquake")
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(
			"DELETE",
			path(1),
			&bytes.Buffer{},
		)
		app.engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Internal server error if something unexpected happens when checking gallery", func(t *testing.T) {
		app := newMockApp()
		app.galleryStore.GetGallery_ = func(input *gallery.GetGalleryInput) (*models.Gallery, error) {
			return nil, errors.New("Nuclear meltdown")
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(
			"DELETE",
			path(1),
			&bytes.Buffer{},
		)
		app.engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
