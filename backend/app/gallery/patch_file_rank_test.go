package app_gallery

import (
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

func TestPatchFileRank(t *testing.T) {
	path := func(fileId string, rank string) string {
		return fmt.Sprintf("/api/files/%s/rank/%s/", fileId, rank)
	}

	t.Run("OK", func(t *testing.T) {
		app := newMockApp()
		app.fileStore.GetFile_ = func(id uint) (*models.File, error) {
			return &models.File{}, nil
		}

		called := false
		app.fileStore.UpdateFileRank_ = func(data *gallery.UpdateFileRankInput) error {
			called = true
			assert.Equal(t, uint(1), data.FileId)
			assert.Equal(t, -2, data.Rank)
			return nil
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(
			"PATCH",
			path("1", "-2"),
			&bytes.Buffer{},
		)
		app.engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.True(t, called)

	})

	t.Run("Internal server error when there is no db", func(t *testing.T) {
		app := newMockApp()
		app.fileStore.GetFile_ = func(id uint) (*models.File, error) {
			return nil, errors.New("Earthquake")
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(
			"PATCH",
			path("1", "-2"),
			&bytes.Buffer{},
		)
		app.engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

	})

	t.Run("Internal server error when something unexpected happens", func(t *testing.T) {
		app := newMockApp()
		app.fileStore.GetFile_ = func(id uint) (*models.File, error) {
			return &models.File{}, nil
		}
		app.fileStore.UpdateFileRank_ = func(data *gallery.UpdateFileRankInput) error {
			return errors.New("Earthquake")
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(
			"PATCH",
			path("1", "-2"),
			&bytes.Buffer{},
		)
		app.engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

	})

	t.Run("Test unauthorized", func(t *testing.T) {
		app := newMockApp()
		app.authContext.RequireUsername_ = func(authHeader string, i auth.IJwt) error {
			return auth.InvalidAuthHeaderError
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(
			"PATCH",
			path("1", "-2"),
			&bytes.Buffer{},
		)
		app.engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Test NOT_FOUND if file doesn't exisits", func(t *testing.T) {
		app := newMockApp()
		app.fileStore.GetFile_ = func(id uint) (*models.File, error) {
			return nil, nil
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(
			"PATCH",
			path("1", "-3"),
			&bytes.Buffer{},
		)
		app.engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Test NOT_FOUND if invalid file id", func(t *testing.T) {
		app := newMockApp()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(
			"PATCH",
			path("invalid", "-2"),
			&bytes.Buffer{},
		)
		app.engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Test NOT_FOUND if invalid rank", func(t *testing.T) {
		app := newMockApp()
		app.fileStore.GetFile_ = func(id uint) (*models.File, error) {
			return &models.File{}, nil
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(
			"PATCH",
			path("2", "asd"),
			&bytes.Buffer{},
		)
		app.engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
