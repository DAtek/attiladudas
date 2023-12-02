package gallery_delete

import (
	"api"
	"api/components/auth"
	"api/components/gallery"
	"bytes"
	"db/models"
	"net/http"

	"github.com/DAtek/gotils"
	"github.com/stretchr/testify/assert"

	"fmt"
	"testing"
)

func TestDeleteGallery(t *testing.T) {
	path := func(id any) string {
		return fmt.Sprintf("/api/gallery/%v/", id)
	}

	t.Run("Delete gallery OK", func(t *testing.T) {
		authCtx := &auth.MockAuthContext{
			RequireUsername_: func(authHeader string) error { return nil },
		}
		galleryId := uint(2)
		galleryStore := &gallery.MockGalleryStore{
			GetGallery_: func(input *gallery.GetGalleryInput) (*models.Gallery, error) {
				return &models.Gallery{Id: galleryId}, nil
			},
			DeleteGallery_: func(gallery *models.Gallery) error {
				assert.Equal(t, galleryId, gallery.Id)
				return nil
			},
		}
		app := api.AppWithMiddlewares(
			PluginDeleteGallery(authCtx, galleryStore),
		)

		req := gotils.ResultOrPanic(http.NewRequest(
			"DELETE",
			path(galleryId),
			&bytes.Buffer{},
		))

		resp := gotils.ResultOrPanic(app.Test(req))

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Test unauthorized", func(t *testing.T) {
		authCtx := &auth.MockAuthContext{
			RequireUsername_: func(authHeader string) error { return auth.InvalidAuthHeaderError },
		}
		app := api.AppWithMiddlewares(
			PluginDeleteGallery(authCtx, &gallery.MockGalleryStore{}),
		)

		req := gotils.ResultOrPanic(http.NewRequest(
			"DELETE",
			path(1),
			&bytes.Buffer{},
		))

		resp := gotils.ResultOrPanic(app.Test(req))

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Not found if gallery not found", func(t *testing.T) {
		authCtx := &auth.MockAuthContext{
			RequireUsername_: func(authHeader string) error { return nil },
		}
		galleryStore := &gallery.MockGalleryStore{
			GetGallery_: func(input *gallery.GetGalleryInput) (*models.Gallery, error) {
				return nil, nil
			},
		}
		app := api.AppWithMiddlewares(
			PluginDeleteGallery(authCtx, galleryStore),
		)

		req := gotils.ResultOrPanic(http.NewRequest(
			"DELETE",
			path(42),
			&bytes.Buffer{},
		))

		resp := gotils.ResultOrPanic(app.Test(req))

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Not found if invalid id", func(t *testing.T) {
		authCtx := &auth.MockAuthContext{
			RequireUsername_: func(authHeader string) error { return nil },
		}
		app := api.AppWithMiddlewares(
			PluginDeleteGallery(authCtx, &gallery.MockGalleryStore{}),
		)

		req := gotils.ResultOrPanic(http.NewRequest(
			"DELETE",
			path("invalid id"),
			&bytes.Buffer{},
		))

		resp := gotils.ResultOrPanic(app.Test(req))

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}
