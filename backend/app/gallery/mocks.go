package app_gallery

import (
	"attiladudas/backend/app/engine"
	"attiladudas/backend/components/auth"
	auth_mocks "attiladudas/backend/components/auth/mocks"
	"attiladudas/backend/components/gallery"
	gallery_mocks "attiladudas/backend/components/gallery/mocks"
	"attiladudas/backend/models"

	"github.com/gin-gonic/gin"
)

type mockApp struct {
	galleryStore *gallery_mocks.MockGalleryStore
	fileStore    *gallery_mocks.MockFileStore
	authContext  *auth_mocks.MockAuthContext
	jwtContext   *auth_mocks.MockJwtContext
	mediaDirName string
	resizer      *gallery_mocks.MockResizer
	engine       *gin.Engine
}

func newMockApp() *mockApp {
	galleryStore := &gallery_mocks.MockGalleryStore{
		GetGallery_:    func(input *gallery.GetGalleryInput) (*models.Gallery, error) { return nil, nil },
		GalleryExists_: func(input *gallery.GetGalleryInput) (bool, error) { return false, nil },
	}
	fileStore := &gallery_mocks.MockFileStore{}
	authContext := &auth_mocks.MockAuthContext{RequireUsername_: func(authHeader string, i auth.IJwt) error { return nil }}

	jwtContext := &auth_mocks.MockJwtContext{
		Encode_: func(c *auth.Claims) (string, error) { return "", nil },
	}

	resizer := &gallery_mocks.MockResizer{}

	mediaDirName := "media"

	appEngine := engine.NewEngine(
		authContext,
		jwtContext,
		&engine.HandlerCollection{
			DeleteFilesHandler:     DeleteFilesHandler(fileStore, galleryStore),
			DeleteGalleryHandler:   DeleteGalleryHandler(galleryStore),
			GetGalleryHandler:      GetGalleryHandler(galleryStore, fileStore),
			GetGalleriesHandler:    GetGalleriesHandler(galleryStore, fileStore, authContext, jwtContext),
			GetResizedImageHandler: GetResizedImageHandler(resizer, mediaDirName),
			PostFilesHandler:       PostFilesHandler(fileStore, galleryStore),
			PatchFileRankHandler:   PatchFileRankHandler(fileStore),
			PostGalleryHandler:     PostGalleryHandler(galleryStore),
			PutGalleryHandler:      PutGalleryHandler(galleryStore),
		},
		mediaDirName,
	)

	return &mockApp{
		authContext:  authContext,
		jwtContext:   jwtContext,
		galleryStore: galleryStore,
		fileStore:    fileStore,
		resizer:      resizer,
		mediaDirName: mediaDirName,
		engine:       appEngine,
	}
}
