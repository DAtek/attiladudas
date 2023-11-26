package room_manager

import (
	"attiladudas/backend/ws"
)

type IRoomManager interface {
	HandleConnection(conn ws.IConn)
}

type actionCollection map[MessageType]action

type action func(manager *roomManager, conn ws.IConn, msg *messageStruct) messageStruct

type roomManager struct {
	roomsByName       map[string]*room
	roomsByConnection map[ws.IConn]*room
	actions           actionCollection
	cleanup_          func(conn ws.IConn)
}
