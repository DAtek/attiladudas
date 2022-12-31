package app_gallery

import (
	"attiladudas/backend/components/gallery"
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetResizedImage(t *testing.T) {
	app := newMockApp()
	path := func(sizePart, filenamePart string) string {
		return fmt.Sprintf("/api/resize/%s/fruits/%s/%s/", app.mediaDirName, sizePart, filenamePart)
	}

	t.Run("Returns resized image", func(t *testing.T) {
		app := newMockApp()
		expectedContent := []byte("image content")
		expectedSize := gallery.Size{Width: 400, Height: 225}

		app.resizer.ResizeImage_ = func(newSize *gallery.Size, directory, filename string) ([]byte, error) {
			if *newSize != expectedSize {
				return nil, errors.New("UNEXPECTED")
			}
			return expectedContent, nil
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(
			"GET",
			path(fmt.Sprintf("%dx%d", expectedSize.Width, expectedSize.Height), "apple.jpg"),
			&bytes.Buffer{},
		)

		app.engine.ServeHTTP(w, req)
		assert.Equal(t, http.StatusFound, w.Code)
	})

	t.Run("Internal server error if resizing fails", func(t *testing.T) {
		app := newMockApp()

		app.resizer.ResizeImage_ = func(newSize *gallery.Size, directory, filename string) ([]byte, error) {
			return nil, errors.New("Unexpected error")
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(
			"GET",
			path("400x225", "apple.jpg"),
			&bytes.Buffer{},
		)

		app.engine.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	invalidSizeScenarios := []string{
		"50x50",
		"asd",
		"x10",
		"50x60x3",
		"225xYOLO",
	}
	for _, size := range invalidSizeScenarios {
		t.Run("Returns not found error if size is invalid", func(t *testing.T) {
			app := newMockApp()

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(
				"GET",
				path(size, "apple.jpg"),
				&bytes.Buffer{},
			)

			app.engine.ServeHTTP(w, req)
			assert.Equal(t, http.StatusNotFound, w.Code)
		})
	}
}
