package app_gallery

import (
	"attiladudas/backend/components"
	"attiladudas/backend/components/auth"
	"attiladudas/backend/components/gallery"
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

func TestDeleteFiles(t *testing.T) {
	path := func(id any) string {
		return fmt.Sprintf("/api/galleries/%v/files/", id)
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

	t.Run("Test OK", func(t *testing.T) {
		app := newMockApp()
		ids := deleteFilesBody{Ids: []uint{1, 2}}
		body, _ := json.Marshal(ids)
		galleryId := uint(2)
		deleteCalled := false

		app.galleryStore.GetGallery_ = func(input *gallery.GetGalleryInput) (*models.Gallery, error) {
			return &models.Gallery{Id: *input.Id}, nil
		}

		app.fileStore.DeleteFiles_ = func(gallery *models.Gallery, fileIds []uint) error {
			deleteCalled = true
			assert.Equal(t, galleryId, gallery.Id)
			assert.Equal(t, ids.Ids, fileIds)
			return nil
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(
			"DELETE",
			path(galleryId),
			bytes.NewBuffer(body),
		)
		app.engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.True(t, deleteCalled)
	})

	t.Run("NOT FOUND if invalid gallery id", func(t *testing.T) {
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

	t.Run("NOT FOUND if gallery doesn't exists", func(t *testing.T) {
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

	t.Run("BAD REQUEST if invalid file ids", func(t *testing.T) {
		app := newMockApp()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(
			"DELETE",
			path(1),
			bytes.NewBuffer([]byte("{\"ids\": [\"a\"]}")),
		)
		app.engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Test INTERNAL SERVER ERROR if something unexpected happens", func(t *testing.T) {
		app := newMockApp()
		ids := deleteFilesBody{Ids: []uint{1, 2}}
		body, _ := json.Marshal(ids)
		app.fileStore.DeleteFiles_ = func(galler *models.Gallery, u2 []uint) error {
			return errors.New("ALIENS ATTACK")
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(
			"DELETE",
			path(1),
			bytes.NewBuffer(body),
		)
		app.engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
