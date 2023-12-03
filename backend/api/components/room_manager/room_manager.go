package room_manager

import ()

type IRoomManager interface {
	HandleConnection(conn IWSConn)
}

type MockRoomManager struct {
	HandleConnection_ func(conn IWSConn)
}

func (m *MockRoomManager) HandleConnection(conn IWSConn) {
	m.HandleConnection_(conn)
}

type actionCollection map[MessageType]action

type action func(manager *roomManager, conn IWSConn, msg *messageStruct) messageStruct

type roomManager struct {
	roomsByName       map[string]*room
	roomsByConnection map[IWSConn]*room
	actions           actionCollection
	cleanup_          func(conn IWSConn)
}
