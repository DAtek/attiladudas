package ws_mocks

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
		ReadChan:  make(chan MockMessage),
		WriteChan: make(chan MockMessage),
	}
}

func (m *MockChanConn) ReadMessage() (messageType int, p []byte, err error) {
	msg := <-m.ReadChan
	return msg.MessageType, msg.Data, msg.Err
}

func (m *MockChanConn) WriteMessage(messageType int, data []byte) error {
	m.WriteChan <- MockMessage{messageType, data, nil}
	return nil
}
