package room_manager

import "attiladudas/backend/ws"

func NewRoomManager() IRoomManager {
	return newRoomManager()
}

func newRoomManager() *roomManager {
	manager := &roomManager{
		roomsByName:       map[string]*room{},
		roomsByConnection: map[ws.IConn]*room{},
		actions:           actions,
	}

	manager.cleanup_ = manager.cleanup
	return manager
}
