package room_manager

func NewRoomManager() IRoomManager {
	return newRoomManager()
}

func newRoomManager() *roomManager {
	manager := &roomManager{
		roomsByName:       map[string]*room{},
		roomsByConnection: map[IWSConn]*room{},
		actions:           actions,
	}

	manager.cleanup_ = manager.cleanup
	return manager
}
