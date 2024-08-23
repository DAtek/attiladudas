package room_manager

type MockChanConn struct {
	ReadChan  chan MockMessage
	WriteChan chan MockMessage
}

type MockMessage struct {
	MessageType int
	Data        []byte
	Err         error
}

func NewMockChanConn() *MockChanConn {
	return &MockChanConn{
		ReadChan:  make(chan MockMessage, 1),
		WriteChan: make(chan MockMessage, 1),
	}
}

func (m *MockChanConn) ReadMessage() (int, []byte, error) {
	msg := <-m.ReadChan
	return msg.MessageType, msg.Data, msg.Err
}

func (m *MockChanConn) WriteMessage(messageType int, data []byte) error {
	m.WriteChan <- MockMessage{messageType, data, nil}
	return nil
}

var _ IWSConn = NewMockChanConn()
