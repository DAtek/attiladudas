package gallery_post

import (
	"api"
	"api/components/auth"
	"api/components/gallery"
	"api/helpers"
	"bytes"
	"db/models"
	"encoding/json"
	"fibertools"
	"io"
	"net/http"
	"testing"

	"github.com/DAtek/gotils"
	"github.com/stretchr/testify/assert"
)

func TestPostGallery(t *testing.T) {
	path := "/api/gallery/"
	validInput := CreateUpdateGalleryBody{
		Title:       "Gallery1",
		Slug:        "gallery-1",
		Description: "Desc1",
		Date:        "2020-01-01",
		Active:      true,
	}

	t.Run("Creates gallery", func(t *testing.T) {
		createCalled := false
		date := gotils.ResultOrPanic(helpers.DateFromISO8601(validInput.Date))
		galleryStore := &gallery.MockGalleryStore{
			GalleryExists_: func(input *gallery.GetGalleryInput) (bool, error) {
				return false, nil
			},
			CreateGallery_: func(i *gallery.CreateUpdateGalleryInput) (*models.Gallery, error) {
				createCalled = true
				wantedInput := &gallery.CreateUpdateGalleryInput{
					Title:       validInput.Title,
					Slug:        validInput.Slug,
					Description: validInput.Description,
					Date:        date,
					Active:      validInput.Active,
				}
				assert.Equal(t, wantedInput, i)
				return &models.Gallery{Id: 1}, nil
			},
		}

		authCtx := &auth.MockAuthContext{
			RequireUsername_: func(authHeader string) error {
				return nil
			},
		}
		app := api.AppWithMiddlewares(
			PluginPostGallery(authCtx, galleryStore),
		)

		body := gotils.ResultOrPanic(json.Marshal(validInput))

		req := gotils.ResultOrPanic(http.NewRequest(
			"POST",
			path,
			bytes.NewBuffer(body),
		))
		req.Header.Add("Content-Type", "application/json")
		resp := gotils.ResultOrPanic(app.Test(req))
		respData := &bytes.Buffer{}
		io.Copy(respData, resp.Body)
		result := &createGalleryResponse{}
		json.Unmarshal(respData.Bytes(), result)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.Equal(t, uint(1), result.Id)
		assert.True(t, createCalled)
	})

	t.Run("Test unauthorized request", func(t *testing.T) {
		authCtx := &auth.MockAuthContext{
			RequireUsername_: func(authHeader string) error {
				return auth.InvalidAuthHeaderError
			},
		}
		app := api.AppWithMiddlewares(
			PluginPostGallery(authCtx, &gallery.MockGalleryStore{}),
		)

		req := gotils.ResultOrPanic(http.NewRequest(
			"POST",
			path,
			&bytes.Buffer{},
		))

		resp := gotils.ResultOrPanic(app.Test(req))

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	badDate := gotils.ResultOrPanic(json.Marshal(CreateUpdateGalleryBody{
		Title:       "Gallery1",
		Description: "Desc1",
		Date:        "asd",
		Active:      true,
	}))

	badRequestBodies := [][]byte{
		[]byte(""),
		badDate,
	}
	for _, body := range badRequestBodies {
		t.Run("Bad request if body is invalid", func(t *testing.T) {
			authCtx := &auth.MockAuthContext{
				RequireUsername_: func(authHeader string) error {
					return nil
				},
			}
			galleryStore := &gallery.MockGalleryStore{
				GalleryExists_: func(input *gallery.GetGalleryInput) (bool, error) {
					return false, nil
				},
			}
			app := api.AppWithMiddlewares(
				PluginPostGallery(authCtx, galleryStore),
			)

			req := gotils.ResultOrPanic(http.NewRequest(
				"POST",
				path,
				bytes.NewBuffer(body),
			))
			req.Header.Add("Content-Type", "application/json")
			resp := gotils.ResultOrPanic(app.Test(req))

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	}

	t.Run("Bad request if gallery with slug exists", func(t *testing.T) {
		galleryStore := &gallery.MockGalleryStore{
			GalleryExists_: func(input *gallery.GetGalleryInput) (bool, error) {
				return true, nil
			},
		}

		authCtx := &auth.MockAuthContext{
			RequireUsername_: func(authHeader string) error {
				return nil
			},
		}
		app := api.AppWithMiddlewares(
			PluginPostGallery(authCtx, galleryStore),
		)

		body := gotils.ResultOrPanic(json.Marshal(validInput))
		req := gotils.ResultOrPanic(http.NewRequest(
			"POST",
			path,
			bytes.NewBuffer(body),
		))
		req.Header.Add("Content-Type", "application/json")
		resp := gotils.ResultOrPanic(app.Test(req))

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Bad request if gallery with same title already exists", func(t *testing.T) {
		galleryStore := &gallery.MockGalleryStore{
			GalleryExists_: func(input *gallery.GetGalleryInput) (bool, error) {
				if input.Title != nil {
					return *input.Title == validInput.Title, nil
				}

				return false, nil
			},
		}

		authCtx := &auth.MockAuthContext{
			RequireUsername_: func(authHeader string) error {
				return nil
			},
		}
		app := api.AppWithMiddlewares(
			PluginPostGallery(authCtx, galleryStore),
		)

		body := gotils.ResultOrPanic(json.Marshal(validInput))
		req := gotils.ResultOrPanic(http.NewRequest(
			"POST",
			path,
			bytes.NewBuffer(body),
		))
		req.Header.Add("Content-Type", "application/json")
		resp := gotils.ResultOrPanic(app.Test(req))

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		respData := &bytes.Buffer{}
		io.Copy(respData, resp.Body)
		err := &fibertools.JsonErrorCollection{}
		json.Unmarshal(respData.Bytes(), err)

		assert.Equal(t, 1, len(err.Errors))
		assert.Equal(t, api.TypeErrorAlreadyExists, err.Errors[0].Type)
		assert.Equal(t, "title", err.Errors[0].Location)
	})
}
