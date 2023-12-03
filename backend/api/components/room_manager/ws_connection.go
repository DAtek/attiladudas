package room_manager

type IWSConn interface {
	ReadMessage() (int, []byte, error)
	WriteMessage(messageType int, data []byte) error
}

type MockWSConn struct {
	ReadMessage_  func() (int, []byte, error)
	WriteMessage_ func(messageType int, data []byte) error
}

func (m *MockWSConn) ReadMessage() (int, []byte, error) {
	return m.ReadMessage_()
}

func (m *MockWSConn) WriteMessage(messageType int, data []byte) error {
	return m.WriteMessage_(messageType, data)
}
