package main

import (
	app_auth "attiladudas/backend/app/auth"
	"attiladudas/backend/app/engine"
	app_fiar "attiladudas/backend/app/five_in_a_row"
	app_gallery "attiladudas/backend/app/gallery"
	"attiladudas/backend/components"
	"attiladudas/backend/components/auth"
	"attiladudas/backend/components/gallery"
	"attiladudas/backend/components/room_manager"
	"attiladudas/backend/ws"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

const retrySeconds = 3

type wsUpgrader struct {
	upgrader *websocket.Upgrader
}

func (u *wsUpgrader) Upgrade(
	w http.ResponseWriter,
	r *http.Request,
	responseHeader http.Header,
) (ws.IConn, error) {
	return u.upgrader.Upgrade(w, r, responseHeader)
}

func main() {
	pemPrivateKey := []byte(components.EnvPrivateKey.Load())
	pemPublicKey := []byte(components.EnvPublicKey.Load())

	jwtContext, jwtErr := auth.NewJwtContext(pemPrivateKey, pemPublicKey)
	if jwtErr != nil {
		panic(jwtErr)
	}

	db, dbErr := components.NewDbFromEnv()
	for dbErr != nil {
		fmt.Printf("Can't connect to db. Error: %v\nRetrying in %ds\n", dbErr, retrySeconds)
		time.Sleep(retrySeconds * time.Second)
		db, dbErr = components.NewDbFromEnv()
	}

	userStore := components.NewUserStore(db)
	mediaDir := components.EnvMediaDir.Load()
	galleryStore := gallery.NewGalleryStore(db, mediaDir)
	fileStore := gallery.NewFileStore(db, mediaDir)
	resizer := gallery.NewResizer(mediaDir)

	upgrader := &wsUpgrader{
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}

	roomManager := room_manager.NewRoomManager()
	mediaDirName := fileStore.MediaDirName()

	router := engine.NewEngine(
		&auth.AuthContext{},
		jwtContext,
		&engine.HandlerCollection{
			DeleteFilesHandler:     app_gallery.DeleteFilesHandler(fileStore, galleryStore),
			DeleteGalleryHandler:   app_gallery.DeleteGalleryHandler(galleryStore),
			FiveInARowHandler:      app_fiar.FiveInARowHandler(upgrader, roomManager),
			GetGalleryHandler:      app_gallery.GetGalleryHandler(galleryStore, fileStore),
			GetGalleriesHandler:    app_gallery.GetGalleriesHandler(galleryStore, fileStore, &auth.AuthContext{}, jwtContext),
			GetResizedImageHandler: app_gallery.GetResizedImageHandler(resizer, mediaDirName),
			PatchFileRankHandler:   app_gallery.PatchFileRankHandler(fileStore),
			PostFilesHandler:       app_gallery.PostFilesHandler(fileStore, galleryStore),
			PostGalleryHandler:     app_gallery.PostGalleryHandler(galleryStore),
			PostTokenHandler:       app_auth.PostTokenHandler(userStore, jwtContext),
			PutGalleryHandler:      app_gallery.PutGalleryHandler(galleryStore),
		},
		mediaDirName,
	)

	socketFile := components.EnvSocketPath.Load()
	if socketFile == "" {
		panic("Variable 'API_SOCKET_PATH' not set.")
	}
	os.Remove(socketFile)
	router.RunUnix(socketFile)
}
