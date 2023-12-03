package main

import (
	"api"
	"api/components/auth"
	"api/components/gallery"
	"api/components/room_manager"
	"api/handlers/file_rank_patch"
	"api/handlers/files_delete"
	"api/handlers/files_post"
	"api/handlers/galleries_get"
	"api/handlers/gallery_delete"
	"api/handlers/gallery_get"
	"api/handlers/gallery_post"
	"api/handlers/gallery_put"
	"api/handlers/resize_get"
	"api/handlers/token_post"
	fiar "api/handlers/ws_five_in_a_row"
	"db"
	"flag"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

const retrySeconds = 3

func main() {
	port := flag.String("port", "8000", "-port 8000")
	host := flag.String("host", "127.0.0.1", "-host 127.0.0.1")
	flag.Parse()
	app := createApp()
	app.Listen(*host + ":" + *port)
}

func createApp() *fiber.App {
	pemPrivateKey := []byte(api.EnvPrivateKey.Load())
	pemPublicKey := []byte(api.EnvPublicKey.Load())

	jwtContext, jwtErr := auth.NewJwtContext(pemPrivateKey, pemPublicKey)
	if jwtErr != nil {
		panic(jwtErr)
	}

	authCtx := auth.NewAuthContext(jwtContext)
	session, dbErr := db.NewDbFromEnv()
	for dbErr != nil {
		fmt.Printf("Can't connect to db. Error: %v\nRetrying in %ds\n", dbErr, retrySeconds)
		time.Sleep(retrySeconds * time.Second)
		session, dbErr = db.NewDbFromEnv()
	}

	mediaDir := api.EnvMediaDir.Load()
	galleryStore := gallery.NewGalleryStore(session, mediaDir)
	fileStore := gallery.NewFileStore(session, mediaDir)
	resizer := gallery.NewResizer(mediaDir)
	roomManager := room_manager.NewRoomManager()
	mediaDirName := fileStore.MediaDirName()

	return api.AppWithMiddlewares(
		file_rank_patch.PluginPatchFileRank(fileStore, authCtx),
		files_delete.PluginDeleteFiles(authCtx, galleryStore, fileStore),
		files_post.PluginPostFiles(authCtx, galleryStore, fileStore),
		galleries_get.PluginGetGalleries(galleryStore, fileStore, authCtx),
		gallery_delete.PluginDeleteGallery(authCtx, galleryStore),
		gallery_get.PluginGetGallery(galleryStore, fileStore),
		gallery_post.PluginPostGallery(authCtx, galleryStore),
		gallery_put.PluginPutGallery(authCtx, galleryStore),
		resize_get.PluginResizeImage(mediaDirName, resizer),
		token_post.PluginTokenPost(session, jwtContext),
		fiar.PluginFiveInARow(roomManager),
	)
}
