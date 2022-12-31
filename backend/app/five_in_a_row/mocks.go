package app_fiar

import (
	"attiladudas/backend/app/engine"
	"attiladudas/backend/ws"
	"net/http"

	"github.com/gin-gonic/gin"
)

// App
type mockApp struct {
	engine      *gin.Engine
	wsUpgrader  *mockWsUpgrader
	roomManager *mockRoomManager
}

func newMockApp() *mockApp {
	upgrader := &mockWsUpgrader{}
	roomManager := &mockRoomManager{}

	appEngine := engine.NewEngine(
		nil,
		nil,
		&engine.HandlerCollection{
			FiveInARowHandler: FiveInARowHandler(upgrader, roomManager),
		},
		"",
	)

	return &mockApp{
		engine:      appEngine,
		wsUpgrader:  upgrader,
		roomManager: roomManager,
	}
}

// WSUpgrader
type mockWsUpgrader struct {
	upgrade func(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (ws.IConn, error)
}

func (m *mockWsUpgrader) Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (ws.IConn, error) {
	return m.upgrade(w, r, responseHeader)
}

// RoomManager
type mockRoomManager struct {
	handleConnection func(conn ws.IConn)
}

func (r *mockRoomManager) HandleConnection(conn ws.IConn) {
	r.handleConnection(conn)
}
