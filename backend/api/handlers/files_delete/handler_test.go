package files_delete

import (
	"api"
	"api/components/auth"
	"api/components/gallery"
	"bytes"
	"db/models"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/DAtek/gotils"
	"github.com/stretchr/testify/assert"
)

func TestDeleteFiles(t *testing.T) {
	path := func(id any) string {
		return fmt.Sprintf("/api/gallery/%v/files/", id)
	}

	t.Run("Test OK", func(t *testing.T) {
		galleryStore := &gallery.MockGalleryStore{
			GetGallery_: func(input *gallery.GetGalleryInput) (*models.Gallery, error) {
				return &models.Gallery{Id: *input.Id}, nil
			},
		}
		deleteCalled := false
		galleryId := uint(42)
		reqData := deleteFilesBody{Ids: []uint{1, 2}}

		fileStore := &gallery.MockFileStore{
			DeleteFiles_: func(gallery *models.Gallery, fileIds []uint) error {
				deleteCalled = true
				assert.Equal(t, gallery.Id, galleryId)
				assert.Equal(t, reqData.Ids, fileIds)
				return nil
			},
		}

		authCtx := auth.MockAuthContext{
			RequireUsername_: func(authHeader string) error { return nil },
		}

		app := api.AppWithMiddlewares(
			PluginDeleteFiles(&authCtx, galleryStore, fileStore),
		)

		body := gotils.ResultOrPanic(json.Marshal(reqData))

		req := gotils.ResultOrPanic(http.NewRequest(
			"DELETE",
			path(galleryId),
			bytes.NewBuffer(body),
		))
		req.Header.Add("Content-Type", "application/json")

		resp := gotils.ResultOrPanic(app.Test(req))

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.True(t, deleteCalled)
	})

	t.Run("Test unauthorized", func(t *testing.T) {
		authCtx := auth.MockAuthContext{
			RequireUsername_: func(authHeader string) error { return auth.InvalidAuthHeaderError },
		}

		app := api.AppWithMiddlewares(
			PluginDeleteFiles(&authCtx, &gallery.MockGalleryStore{}, &gallery.MockFileStore{}),
		)

		req := gotils.ResultOrPanic(http.NewRequest(
			"DELETE",
			path(42),
			&bytes.Buffer{},
		))

		resp := gotils.ResultOrPanic(app.Test(req))

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("NOT FOUND if invalid gallery id", func(t *testing.T) {
		authCtx := auth.MockAuthContext{
			RequireUsername_: func(authHeader string) error { return nil },
		}

		app := api.AppWithMiddlewares(
			PluginDeleteFiles(&authCtx, &gallery.MockGalleryStore{}, &gallery.MockFileStore{}),
		)

		req := gotils.ResultOrPanic(http.NewRequest(
			"DELETE",
			path("invalid id LOL"),
			&bytes.Buffer{},
		))

		resp := gotils.ResultOrPanic(app.Test(req))

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("NOT FOUND if gallery doesn't exists", func(t *testing.T) {
		galleryStore := &gallery.MockGalleryStore{
			GetGallery_: func(input *gallery.GetGalleryInput) (*models.Gallery, error) { return nil, nil },
		}

		authCtx := auth.MockAuthContext{
			RequireUsername_: func(authHeader string) error { return nil },
		}

		app := api.AppWithMiddlewares(
			PluginDeleteFiles(&authCtx, galleryStore, &gallery.MockFileStore{}),
		)

		req := gotils.ResultOrPanic(http.NewRequest(
			"DELETE",
			path(33),
			&bytes.Buffer{},
		))

		resp := gotils.ResultOrPanic(app.Test(req))

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("BAD REQUEST if invalid file ids", func(t *testing.T) {
		authCtx := auth.MockAuthContext{
			RequireUsername_: func(authHeader string) error { return nil },
		}

		galleryStore := &gallery.MockGalleryStore{
			GetGallery_: func(input *gallery.GetGalleryInput) (*models.Gallery, error) {
				return &models.Gallery{}, nil
			},
		}

		app := api.AppWithMiddlewares(
			PluginDeleteFiles(&authCtx, galleryStore, &gallery.MockFileStore{}),
		)

		req := gotils.ResultOrPanic(http.NewRequest(
			"DELETE",
			path(1),
			bytes.NewBuffer([]byte("{\"ids\": [\"a\"]}")),
		))
		req.Header.Add("Content-Type", "application/json")

		resp := gotils.ResultOrPanic(app.Test(req))

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Internal server error if deletion fails", func(t *testing.T) {
		galleryStore := &gallery.MockGalleryStore{
			GetGallery_: func(input *gallery.GetGalleryInput) (*models.Gallery, error) {
				return &models.Gallery{Id: *input.Id}, nil
			},
		}
		deleteCalled := false
		galleryId := uint(42)
		reqData := deleteFilesBody{Ids: []uint{1, 2}}

		fileStore := &gallery.MockFileStore{
			DeleteFiles_: func(gallery *models.Gallery, fileIds []uint) error {
				deleteCalled = true
				return errors.New("unexpected")
			},
		}

		authCtx := auth.MockAuthContext{
			RequireUsername_: func(authHeader string) error { return nil },
		}

		app := api.AppWithMiddlewares(
			PluginDeleteFiles(&authCtx, galleryStore, fileStore),
		)

		body := gotils.ResultOrPanic(json.Marshal(reqData))

		req := gotils.ResultOrPanic(http.NewRequest(
			"DELETE",
			path(galleryId),
			bytes.NewBuffer(body),
		))
		req.Header.Add("Content-Type", "application/json")

		resp := gotils.ResultOrPanic(app.Test(req))

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		assert.True(t, deleteCalled)
	})
}
