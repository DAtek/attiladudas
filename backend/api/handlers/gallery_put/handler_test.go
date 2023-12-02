package gallery_put

import (
	"api"
	"api/components/auth"
	"api/components/gallery"
	"api/handlers/shared"
	"api/helpers"
	"bytes"
	"db/models"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/DAtek/gotils"
	"github.com/stretchr/testify/assert"
)

func TestPutGallery(t *testing.T) {
	path := func(id any) string {
		return fmt.Sprintf("/api/gallery/%v/", id)
	}

	validInput := shared.CreateUpdateGalleryBody{
		Title:       "Gallery1",
		Slug:        "gallery-1",
		Description: "Desc1",
		Date:        "2020-01-01",
		Active:      true,
	}
	validId := uint(2)
	validBody := gotils.ResultOrPanic(json.Marshal(validInput))
	t.Run("Update gallery", func(t *testing.T) {
		authCtx := &auth.MockAuthContext{
			RequireUsername_: func(authHeader string) error {
				return nil
			},
		}

		updateCalled := false
		galleryStore := &gallery.MockGalleryStore{
			GalleryExists_: func(input *gallery.GetGalleryInput) (bool, error) { return true, nil },
			GetGallery_: func(input *gallery.GetGalleryInput) (*models.Gallery, error) {
				return &models.Gallery{
					Id:    *input.Id,
					Title: validInput.Title,
					Slug:  validInput.Slug,
				}, nil
			},
			UpdateGallery_: func(galleryId uint, input *gallery.CreateUpdateGalleryInput) error {
				updateCalled = true
				assert.Equal(t, validId, galleryId)
				date := gotils.ResultOrPanic(helpers.DateFromISO8601(validInput.Date))
				wantedInput := &gallery.CreateUpdateGalleryInput{
					Title:       validInput.Title,
					Slug:        validInput.Slug,
					Description: validInput.Description,
					Date:        date,
					Active:      validInput.Active,
				}
				assert.Equal(t, wantedInput, input)
				return nil
			},
		}

		app := api.AppWithMiddlewares(
			PluginPutGallery(authCtx, galleryStore),
		)

		req := gotils.ResultOrPanic(http.NewRequest(
			"PUT",
			path(validId),
			bytes.NewBuffer([]byte(validBody)),
		))
		req.Header.Add("content-type", "application/json")

		resp := gotils.ResultOrPanic(app.Test(req))

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.True(t, updateCalled)
	})

	t.Run("Test unauthorized", func(t *testing.T) {
		authCtx := &auth.MockAuthContext{
			RequireUsername_: func(authHeader string) error {
				return auth.InvalidAuthHeaderError
			},
		}

		app := api.AppWithMiddlewares(
			PluginPutGallery(authCtx, &gallery.MockGalleryStore{}),
		)

		req := gotils.ResultOrPanic(http.NewRequest(
			"PUT",
			path(validId),
			&bytes.Buffer{},
		))

		resp := gotils.ResultOrPanic(app.Test(req))

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	badDate := gotils.ResultOrPanic(json.Marshal(shared.CreateUpdateGalleryBody{
		Title:       "Gallery1",
		Description: "Desc1",
		Date:        "1",
		Active:      true,
	}))
	badRequestscenarios := [][]byte{
		[]byte(""),
		badDate,
	}
	for _, body := range badRequestscenarios {
		t.Run("Bad request if body is invalid", func(t *testing.T) {
			authCtx := &auth.MockAuthContext{
				RequireUsername_: func(authHeader string) error {
					return nil
				},
			}

			galleryStore := &gallery.MockGalleryStore{
				GetGallery_: func(input *gallery.GetGalleryInput) (*models.Gallery, error) {
					return &models.Gallery{}, nil
				},
				GalleryExists_: func(input *gallery.GetGalleryInput) (bool, error) { return true, nil },
			}

			app := api.AppWithMiddlewares(
				PluginPutGallery(authCtx, galleryStore),
			)
			req := gotils.ResultOrPanic(http.NewRequest(
				"PUT",
				path(validId),
				bytes.NewBuffer(body),
			))
			req.Header.Add("content-type", "application/json")

			resp := gotils.ResultOrPanic(app.Test(req))

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	}

	t.Run("Bad request if gallery with same title already exists", func(t *testing.T) {
		authCtx := &auth.MockAuthContext{
			RequireUsername_: func(authHeader string) error {
				return nil
			},
		}

		galleryStore := &gallery.MockGalleryStore{
			GetGallery_: func(input *gallery.GetGalleryInput) (*models.Gallery, error) {
				return &models.Gallery{}, nil
			},
			GalleryExists_: func(input *gallery.GetGalleryInput) (bool, error) { return true, nil },
		}

		app := api.AppWithMiddlewares(
			PluginPutGallery(authCtx, galleryStore),
		)

		body := gotils.ResultOrPanic(json.Marshal(validInput))
		req := gotils.ResultOrPanic(http.NewRequest(
			"PUT",
			path(validId),
			bytes.NewBuffer(body),
		))
		req.Header.Add("content-type", "application/json")

		resp := gotils.ResultOrPanic(app.Test(req))

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Not found if id is invalid", func(t *testing.T) {
		authCtx := &auth.MockAuthContext{
			RequireUsername_: func(authHeader string) error {
				return nil
			},
		}

		app := api.AppWithMiddlewares(
			PluginPutGallery(authCtx, &gallery.MockGalleryStore{}),
		)
		req := gotils.ResultOrPanic(http.NewRequest(
			"PUT",
			path("asd"),
			&bytes.Buffer{},
		))

		resp := gotils.ResultOrPanic(app.Test(req))

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Not found if gallery doesn't exist", func(t *testing.T) {
		authCtx := &auth.MockAuthContext{
			RequireUsername_: func(authHeader string) error {
				return nil
			},
		}

		galleryStore := &gallery.MockGalleryStore{
			GetGallery_: func(input *gallery.GetGalleryInput) (*models.Gallery, error) { return nil, nil },
		}

		app := api.AppWithMiddlewares(
			PluginPutGallery(authCtx, galleryStore),
		)
		req := gotils.ResultOrPanic(http.NewRequest(
			"PUT",
			path(1),
			&bytes.Buffer{},
		))

		resp := gotils.ResultOrPanic(app.Test(req))

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}
