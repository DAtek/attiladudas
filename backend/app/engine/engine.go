package engine

import (
	"attiladudas/backend/components/auth"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HandlerCollection struct {
	DeleteFilesHandler     gin.HandlerFunc
	DeleteGalleryHandler   gin.HandlerFunc
	FiveInARowHandler      gin.HandlerFunc
	GetGalleryHandler      gin.HandlerFunc
	GetGalleriesHandler    gin.HandlerFunc
	GetResizedImageHandler gin.HandlerFunc
	PatchFileRankHandler   gin.HandlerFunc
	PostFilesHandler       gin.HandlerFunc
	PostGalleryHandler     gin.HandlerFunc
	PostTokenHandler       gin.HandlerFunc
	PutGalleryHandler      gin.HandlerFunc
}

func NewEngine(
	authContext auth.IAuthorization,
	jwtContext auth.IJwt,
	handlerCollection *HandlerCollection,
	mediaDirName string,
) *gin.Engine {
	engine := gin.Default()
	engine.SetTrustedProxies(nil)

	requireUsername := func(ctx *gin.Context) {
		authHeader := ctx.Request.Header.Get("authorization")
		if err := authContext.RequireUsername(authHeader, jwtContext); err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}
	}

	// engine.PATCH("/api/files/:id/rank/:rank/", requireUsername, handlerCollection.PatchFileRankHandler)
	// engine.GET("/api/gallery/:slug/", handlerCollection.GetGalleryHandler)
	// engine.POST("/api/gallery/", requireUsername, handlerCollection.PostGalleryHandler)
	// engine.GET("/api/galleries/", handlerCollection.GetGalleriesHandler)
	// engine.PUT("/api/gallery/:id/", requireUsername, handlerCollection.PutGalleryHandler)
	// engine.DELETE("/api/gallery/:id/", requireUsername, handlerCollection.DeleteGalleryHandler)
	// engine.POST("/api/gallery/:id/files/", requireUsername, handlerCollection.PostFilesHandler)
	// engine.DELETE("/api/gallery/:id/files/", requireUsername, handlerCollection.DeleteFilesHandler)
	// engine.POST("/api/token/", handlerCollection.PostTokenHandler)
	engine.GET(
		fmt.Sprintf("/api/resize/%s/:directory/:size/:filename/", mediaDirName),
		handlerCollection.GetResizedImageHandler,
	)
	// engine.GET("/ws/five-in-a-row/", handlerCollection.FiveInARowHandler)

	return engine
}
