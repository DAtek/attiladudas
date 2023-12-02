package files_post

import (
	"api"
	"api/components/auth"
	"api/components/gallery"
	"bytes"
	"db/models"
	"encoding/json"
	"fibertools"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"testing"

	"github.com/DAtek/gotils"
	"github.com/stretchr/testify/assert"
)

func TestPostFiles(t *testing.T) {
	path := func(id any) string {
		return fmt.Sprintf("/api/gallery/%v/files/", id)
	}

	t.Run("Test upload files", func(t *testing.T) {
		authCtx := &auth.MockAuthContext{
			RequireUsername_: func(authHeader string) error { return nil },
		}

		galleryId := uint(42)
		bodyWithFiles := newMultipartBody()

		fileStore := &gallery.MockFileStore{
			AddFiles_: func(g *models.Gallery, fi []*gallery.FileInput) error {
				assert.Equal(t, galleryId, g.Id)
				for i := range bodyWithFiles.files {
					// Wanted content is being saved
					assert.Equal(t, bodyWithFiles.files[i], fi[i])
				}
				return nil
			},
		}

		galleryStore := &gallery.MockGalleryStore{
			GetGallery_: func(input *gallery.GetGalleryInput) (*models.Gallery, error) {
				return &models.Gallery{Id: *input.Id}, nil
			},
		}

		app := api.AppWithMiddlewares(
			PluginPostFiles(authCtx, galleryStore, fileStore),
		)

		req := gotils.ResultOrPanic(http.NewRequest(
			"POST",
			path(galleryId),
			bodyWithFiles.body,
		))

		req.Header.Set("content-type", bodyWithFiles.multipartWriter.FormDataContentType())

		resp := gotils.ResultOrPanic(app.Test(req))
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	t.Run("Test unauthorized", func(t *testing.T) {
		authCtx := &auth.MockAuthContext{
			RequireUsername_: func(authHeader string) error { return auth.InvalidClaimsError },
		}

		app := api.AppWithMiddlewares(
			PluginPostFiles(authCtx, &gallery.MockGalleryStore{}, &gallery.MockFileStore{}),
		)

		req := gotils.ResultOrPanic(http.NewRequest(
			"POST",
			path(42),
			&bytes.Buffer{},
		))

		resp := gotils.ResultOrPanic(app.Test(req))
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Bad request if filename already exists", func(t *testing.T) {
		authCtx := &auth.MockAuthContext{
			RequireUsername_: func(authHeader string) error { return nil },
		}

		galleryId := uint(42)
		bodyWithFiles := newMultipartBody("apple.jpg", "alien.jpg")

		galleryStore := &gallery.MockGalleryStore{
			GetGallery_: func(input *gallery.GetGalleryInput) (*models.Gallery, error) {
				return &models.Gallery{
					Id: *input.Id,
					Files: []*models.File{
						{Filename: "alien.jpg"},
					},
				}, nil
			},
		}

		app := api.AppWithMiddlewares(
			PluginPostFiles(authCtx, galleryStore, &gallery.MockFileStore{}),
		)

		req := gotils.ResultOrPanic(http.NewRequest(
			"POST",
			path(galleryId),
			bodyWithFiles.body,
		))

		req.Header.Set("content-type", bodyWithFiles.multipartWriter.FormDataContentType())

		resp := gotils.ResultOrPanic(app.Test(req))
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		respData := &bytes.Buffer{}
		io.Copy(respData, resp.Body)
		err := &fibertools.JsonErrorCollection{}
		json.Unmarshal(respData.Bytes(), err)
		assert.Equal(t, "files.1.filename", err.Errors[0].Location)
	})

	t.Run("Bad request if content type is wrong", func(t *testing.T) {
		authCtx := &auth.MockAuthContext{
			RequireUsername_: func(authHeader string) error { return nil },
		}

		galleryId := uint(42)
		bodyWithFiles := newMultipartBody("apple.jpg", "alien.jpg")

		galleryStore := &gallery.MockGalleryStore{
			GetGallery_: func(input *gallery.GetGalleryInput) (*models.Gallery, error) {
				return &models.Gallery{
					Id: *input.Id,
					Files: []*models.File{
						{Filename: "alien.jpg"},
					},
				}, nil
			},
		}

		app := api.AppWithMiddlewares(
			PluginPostFiles(authCtx, galleryStore, &gallery.MockFileStore{}),
		)

		req := gotils.ResultOrPanic(http.NewRequest(
			"POST",
			path(galleryId),
			bodyWithFiles.body,
		))

		resp := gotils.ResultOrPanic(app.Test(req))
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Not found if invalid id", func(t *testing.T) {
		authCtx := &auth.MockAuthContext{
			RequireUsername_: func(authHeader string) error { return nil },
		}

		app := api.AppWithMiddlewares(
			PluginPostFiles(authCtx, &gallery.MockGalleryStore{}, &gallery.MockFileStore{}),
		)

		req := gotils.ResultOrPanic(http.NewRequest(
			"POST",
			path("invalid id"),
			&bytes.Buffer{},
		))

		resp := gotils.ResultOrPanic(app.Test(req))
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Not found if gallery doesn't exist", func(t *testing.T) {
		authCtx := &auth.MockAuthContext{
			RequireUsername_: func(authHeader string) error { return nil },
		}

		galleryStore := &gallery.MockGalleryStore{
			GetGallery_: func(input *gallery.GetGalleryInput) (*models.Gallery, error) { return nil, nil },
		}

		app := api.AppWithMiddlewares(
			PluginPostFiles(authCtx, galleryStore, &gallery.MockFileStore{}),
		)

		req := gotils.ResultOrPanic(http.NewRequest(
			"POST",
			path(42),
			&bytes.Buffer{},
		))

		resp := gotils.ResultOrPanic(app.Test(req))
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

type multipartBody struct {
	files           []*gallery.FileInput
	multipartWriter *multipart.Writer
	body            *bytes.Buffer
}

func newMultipartBody(filenames ...interface{}) *multipartBody {
	files := []*gallery.FileInput{}

	if filenames == nil {
		files = []*gallery.FileInput{
			{Content: []byte("batman"), Filename: "batman.jpg"},
			{Content: []byte("alien"), Filename: "alien.jpg"},
		}
	} else {
		for _, item := range filenames {
			filename, ok := item.(string)
			if !ok {
				panic("filename isn't string")
			}
			files = append(
				files,
				&gallery.FileInput{Content: []byte(filename), Filename: filename})
		}
	}

	bodyWithFiles := &bytes.Buffer{}
	multipartWriter := multipart.NewWriter(bodyWithFiles)
	for _, file := range files {
		part, _ := multipartWriter.CreateFormFile("files", file.Filename)
		part.Write(file.Content)
	}
	multipartWriter.Close()

	return &multipartBody{
		files:           files,
		multipartWriter: multipartWriter,
		body:            bodyWithFiles,
	}
}
