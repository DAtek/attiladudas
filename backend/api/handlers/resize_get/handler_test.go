package resize_get

import (
	"api"
	"api/components/gallery"
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/DAtek/gotils"
	"github.com/stretchr/testify/assert"
)

func TestGetResizedImage(t *testing.T) {
	mediaDir := api.EnvMediaDir.Load()
	path := func(sizePart, filenamePart string) string {
		return fmt.Sprintf("/api/resize/%s/fruits/%s/%s/", mediaDir, sizePart, filenamePart)
	}

	t.Run("Returns resized image", func(t *testing.T) {
		expectedSize := gallery.Size{Width: 400, Height: 225}
		expectedContent := []byte("image content")

		resizer := &gallery.MockResizer{
			ResizeImage_: func(newSize *gallery.Size, directory, filename string) ([]byte, error) {
				if *newSize != expectedSize {
					return nil, errors.New("UNEXPECTED_ERROR")
				}
				return expectedContent, nil
			},
		}

		app := api.AppWithMiddlewares(
			PluginResizeImage(mediaDir, resizer),
		)

		req := gotils.ResultOrPanic(http.NewRequest(
			"GET",
			path(fmt.Sprintf("%dx%d", expectedSize.Width, expectedSize.Height), "apple.jpg"),
			&bytes.Buffer{},
		))

		resp := gotils.ResultOrPanic(app.Test(req))

		assert.Equal(t, http.StatusFound, resp.StatusCode)
	})

	t.Run("Internal server error if resizing fails", func(t *testing.T) {
		resizer := &gallery.MockResizer{
			ResizeImage_: func(newSize *gallery.Size, directory, filename string) ([]byte, error) {
				return nil, errors.New("Unexpected error")
			},
		}

		app := api.AppWithMiddlewares(
			PluginResizeImage(mediaDir, resizer),
		)

		req := gotils.ResultOrPanic(http.NewRequest(
			"GET",
			path("400x225", "apple.jpg"),
			&bytes.Buffer{},
		))

		resp := gotils.ResultOrPanic(app.Test(req))

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
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
			app := api.AppWithMiddlewares(
				PluginResizeImage(mediaDir, &gallery.MockResizer{}),
			)

			req := gotils.ResultOrPanic(http.NewRequest(
				"GET",
				path(size, "apple.jpg"),
				&bytes.Buffer{},
			))

			resp := gotils.ResultOrPanic(app.Test(req))
			assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		})
	}
}
