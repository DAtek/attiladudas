package file_rank_patch

import (
	"api"
	"api/components/auth"
	"api/components/gallery"
	"bytes"
	"db/models"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/DAtek/gotils"
	"github.com/stretchr/testify/assert"
)

func TestPatchFileRank(t *testing.T) {
	getPath := func(fileId string, rank string) string {
		return fmt.Sprintf("/api/file/%s/rank/%s/", fileId, rank)
	}

	t.Run("OK", func(t *testing.T) {
		file := &models.File{Id: 1}
		fileStore := &gallery.MockFileStore{
			GetFile_: func(id uint) (*models.File, error) {
				if id == file.Id {
					return file, nil
				}
				return nil, errors.New("UNEXPECTED")
			},
			UpdateFileRank_: func(data *gallery.UpdateFileRankInput) error {
				file.Rank = data.Rank
				return nil
			},
		}

		authCtx := &auth.MockAuthContext{
			RequireUsername_: func(authHeader string) error {
				return nil
			},
		}

		app := api.AppWithMiddlewares(
			PluginPatchFileRank(fileStore, authCtx),
		)

		req := gotils.ResultOrPanic(http.NewRequest(
			"PATCH",
			getPath("1", "-2"),
			&bytes.Buffer{},
		))

		resp := gotils.ResultOrPanic(app.Test(req))

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, -2, file.Rank)
	})

	t.Run("Internal server error when there is no db", func(t *testing.T) {
		fileStore := &gallery.MockFileStore{
			GetFile_: func(id uint) (*models.File, error) {
				return nil, errors.New("UNEXPECTED")
			},
		}
		authCtx := &auth.MockAuthContext{
			RequireUsername_: func(authHeader string) error {
				return nil
			},
		}

		app := api.AppWithMiddlewares(
			PluginPatchFileRank(fileStore, authCtx),
		)

		req := gotils.ResultOrPanic(http.NewRequest(
			"PATCH",
			getPath("1", "-2"),
			&bytes.Buffer{},
		))

		resp := gotils.ResultOrPanic(app.Test(req))

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("Test unauthorized", func(t *testing.T) {
		authCtx := &auth.MockAuthContext{
			RequireUsername_: func(authHeader string) error {
				return auth.InvalidAuthHeaderError
			},
		}

		app := api.AppWithMiddlewares(
			PluginPatchFileRank(&gallery.MockFileStore{}, authCtx),
		)

		req := gotils.ResultOrPanic(http.NewRequest(
			"PATCH",
			getPath("1", "-2"),
			&bytes.Buffer{},
		))

		resp := gotils.ResultOrPanic(app.Test(req))

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Test NOT_FOUND if file doesn't exisits", func(t *testing.T) {
		fileStore := &gallery.MockFileStore{
			GetFile_: func(id uint) (*models.File, error) {
				return nil, nil
			},
		}
		authCtx := &auth.MockAuthContext{
			RequireUsername_: func(authHeader string) error {
				return nil
			},
		}

		app := api.AppWithMiddlewares(
			PluginPatchFileRank(fileStore, authCtx),
		)

		req := gotils.ResultOrPanic(http.NewRequest(
			"PATCH",
			getPath("1", "-2"),
			&bytes.Buffer{},
		))

		resp := gotils.ResultOrPanic(app.Test(req))

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Test NOT_FOUND if invalid file id", func(t *testing.T) {
		fileStore := &gallery.MockFileStore{
			GetFile_: func(id uint) (*models.File, error) {
				return nil, nil
			},
		}
		authCtx := &auth.MockAuthContext{
			RequireUsername_: func(authHeader string) error {
				return nil
			},
		}

		app := api.AppWithMiddlewares(
			PluginPatchFileRank(fileStore, authCtx),
		)

		req := gotils.ResultOrPanic(http.NewRequest(
			"PATCH",
			getPath("asd", "-2"),
			&bytes.Buffer{},
		))

		resp := gotils.ResultOrPanic(app.Test(req))

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}
