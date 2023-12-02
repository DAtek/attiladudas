package galleries_get

import (
	"api"
	"api/components/auth"
	"api/components/gallery"
	"api/handlers/gallery_get"
	"api/helpers"
	"bytes"
	"db/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/DAtek/gotils"
	"github.com/stretchr/testify/assert"
)

func TestGetGalleries(t *testing.T) {
	path := "/api/galleries/"

	t.Run("Returns galleries", func(t *testing.T) {
		dateString := "2022-01-01"
		date := gotils.ResultOrPanic(helpers.DateFromISO8601(dateString))
		galleries := []*models.Gallery{
			{
				Id:          1,
				Title:       "Gallery1",
				Description: "desc1",
				Date:        date,
				Active:      true,
				Files: []*models.File{
					{Filename: "1.jpg", Rank: 3},
				},
			},
			{
				Id:    2,
				Title: "Gallery2",
				Date:  date,
			},
		}

		wantedResponse := &GalleriesResponse{
			Galleries: []*gallery_get.GalleryResponse{
				{
					Id:          galleries[0].Id,
					Title:       galleries[0].Title,
					Description: galleries[0].Description,
					Active:      galleries[0].Active,
					Date:        dateString,
					Files: []*gallery_get.File{
						{
							Filename: galleries[0].Files[0].Filename,
							Path:     "fake path",
							Rank:     galleries[0].Files[0].Rank,
						},
					},
				},
				{
					Id:    galleries[1].Id,
					Title: galleries[1].Title,
					Files: []*gallery_get.File{},
					Date:  dateString,
				},
			},
			Total: 2,
		}

		galleryStore := &gallery.MockGalleryStore{
			GetGalleries_: func(input *gallery.GetGalleriesInput) (*gallery.PaginatedGalleriesResult, error) {
				return &gallery.PaginatedGalleriesResult{Items: galleries, Total: uint(len(wantedResponse.Galleries))}, nil
			},
		}

		fileStore := &gallery.MockFileStore{
			GetFileDownloadPath_: func(gallery *models.Gallery, filename string) string {
				return wantedResponse.Galleries[0].Files[0].Path
			},
		}

		authCtx := &auth.MockAuthContext{
			RequireUsername_: func(authHeader string) error { return nil },
		}

		app := api.AppWithMiddlewares(
			PluginGetGalleries(galleryStore, fileStore, authCtx),
		)

		wantedInput := &gallery.GetGalleriesInput{
			Page:     2,
			PageSize: 5,
		}

		req, _ := http.NewRequest(
			"GET",
			path,
			&bytes.Buffer{},
		)
		query := req.URL.Query()
		query.Add("page_size", fmt.Sprint(wantedInput.PageSize))
		query.Add("page", fmt.Sprint(wantedInput.Page))
		req.URL.RawQuery = query.Encode()

		resp := gotils.ResultOrPanic(app.Test(req))

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		respData := &bytes.Buffer{}
		io.Copy(respData, resp.Body)

		response := &GalleriesResponse{}
		json.Unmarshal(respData.Bytes(), response)
		assert.Equal(t, wantedResponse, response)
	})

	invalidScenarios := []struct {
		name     string
		pageSize *string
		page     *string
	}{
		{"page size invalid 1", gotils.Pointer("-2"), gotils.Pointer("1")},
		{"page size invalid 2", nil, gotils.Pointer("1")},
		{"page invalid", gotils.Pointer("1"), gotils.Pointer("-1")},
	}

	for _, s := range invalidScenarios {
		t.Run("Returns BAD REQUEST if "+s.name, func(t *testing.T) {
			app := api.AppWithMiddlewares(
				PluginGetGalleries(
					&gallery.MockGalleryStore{},
					&gallery.MockFileStore{},
					&auth.MockAuthContext{
						RequireUsername_: func(authHeader string) error { return nil },
					},
				),
			)

			req, _ := http.NewRequest(
				"GET",
				path,
				&bytes.Buffer{},
			)
			query := req.URL.Query()

			if s.pageSize != nil {
				query.Add("page_size", *s.pageSize)
			}

			if s.page != nil {
				query.Add("page", *s.page)
			}

			req.URL.RawQuery = query.Encode()

			resp := gotils.ResultOrPanic(app.Test(req))

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	}

	t.Run("Returns only active galeries when user is not logged in", func(t *testing.T) {
		activeFilterEbabled := false

		galleryStore := &gallery.MockGalleryStore{
			GetGalleries_: func(input *gallery.GetGalleriesInput) (*gallery.PaginatedGalleriesResult, error) {
				activeFilterEbabled = *input.Active
				return &gallery.PaginatedGalleriesResult{}, nil
			},
		}

		fileStore := &gallery.MockFileStore{}

		authCtx := &auth.MockAuthContext{
			RequireUsername_: func(authHeader string) error {
				return auth.InvalidAuthHeaderError
			},
		}

		app := api.AppWithMiddlewares(
			PluginGetGalleries(galleryStore, fileStore, authCtx),
		)

		req, _ := http.NewRequest(
			"GET",
			path,
			&bytes.Buffer{},
		)
		query := req.URL.Query()
		query.Add("page_size", "1")
		query.Add("page", "1")
		req.URL.RawQuery = query.Encode()

		resp := gotils.ResultOrPanic(app.Test(req))

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.True(t, activeFilterEbabled)
	})
}
