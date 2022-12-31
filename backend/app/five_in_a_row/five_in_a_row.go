package app_fiar

import (
	"attiladudas/backend/components/room_manager"
	"attiladudas/backend/ws"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func FiveInARowHandler(wsupgrader ws.IUpgrader, roomManager room_manager.IRoomManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fiveInARow(ctx, wsupgrader, roomManager)
	}
}

func fiveInARow(ctx *gin.Context, wsupgrader ws.IUpgrader, roomManager room_manager.IRoomManager) {
	conn, err := wsupgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		fmt.Printf("Failed to set websocket upgrade: %v\n", err)
		return
	}

	roomManager.HandleConnection(conn)
}
