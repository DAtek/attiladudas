package app_gallery

import (
	"attiladudas/backend/components"
	"attiladudas/backend/components/auth"
	"attiladudas/backend/components/gallery"
	"attiladudas/backend/helpers"
	"attiladudas/backend/models"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostFiles(t *testing.T) {
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
			"POST",
			path(1),
			&bytes.Buffer{},
		)
		app.engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Bad request if filename already exists", func(t *testing.T) {
		app := newMockApp()
		galleryId := uint(2)
		bodyWithFiles := newMultipartBody("apple.jpg", "alien.jpg")

		app.galleryStore.GetGallery_ = func(input *gallery.GetGalleryInput) (*models.Gallery, error) {
			return &models.Gallery{
				Id: *input.Id,
				Files: []*models.File{
					{Filename: "alien.jpg"},
				},
			}, nil
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(
			"POST",
			path(galleryId),
			bodyWithFiles.body,
		)
		req.Header.Set("content-type", bodyWithFiles.multipartWriter.FormDataContentType())
		app.engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		err := &helpers.JsonErrorCollection{}
		json.Unmarshal(w.Body.Bytes(), err)
		assert.Equal(t, "files.1.filename", *err.Errors[0].Location)
	})

	t.Run("Test upload files", func(t *testing.T) {
		app := newMockApp()
		galleryId := uint(2)
		bodyWithFiles := newMultipartBody()

		app.galleryStore.GetGallery_ = func(input *gallery.GetGalleryInput) (*models.Gallery, error) {
			return &models.Gallery{Id: *input.Id}, nil
		}

		app.fileStore.AddFiles_ = func(g *models.Gallery, fi []*gallery.FileInput) error {
			assert.Equal(t, galleryId, g.Id)
			for i := range bodyWithFiles.files {
				// Wanted content is being saved
				assert.Equal(t, bodyWithFiles.files[i], fi[i])
			}
			return nil
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(
			"POST",
			path(galleryId),
			bodyWithFiles.body,
		)
		req.Header.Set("content-type", bodyWithFiles.multipartWriter.FormDataContentType())
		app.engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("Not found if invalid id", func(t *testing.T) {
		app := newMockApp()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(
			"POST",
			path("asd"),
			&bytes.Buffer{},
		)
		app.engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Not found if gallery doesn't exist", func(t *testing.T) {
		app := newMockApp()
		app.galleryStore.GetGallery_ = func(input *gallery.GetGalleryInput) (*models.Gallery, error) {
			return nil, components.NotFoundError
		}
		bodyWithFiles := newMultipartBody()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(
			"POST",
			path(1),
			bodyWithFiles.body,
		)
		req.Header.Set("content-type", bodyWithFiles.multipartWriter.FormDataContentType())
		app.engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Internal server error if something bad happens", func(t *testing.T) {
		app := newMockApp()
		app.fileStore.AddFiles_ = func(g *models.Gallery, fi []*gallery.FileInput) error {
			return errors.New("!!!FIRE!!!")
		}
		bodyWithFiles := newMultipartBody()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(
			"POST",
			path(1),
			bodyWithFiles.body,
		)
		req.Header.Set("content-type", bodyWithFiles.multipartWriter.FormDataContentType())
		app.engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
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
