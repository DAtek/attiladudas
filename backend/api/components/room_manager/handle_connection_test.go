package room_manager

import (
	"encoding/json"
	"errors"
	"sync"
	"testing"

	"github.com/DAtek/gotils"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestHandleConnection(t *testing.T) {
	t.Run("Returns when message type is ws close", func(t *testing.T) {
		conn := NewMockChanConn()
		timeout := gotils.NewTimeoutMs(100)
		go func() { panic(<-timeout.ErrorCh) }()
		defer timeout.Cancel()
		roomManager := newRoomManager()
		cleanupCalled := false
		roomManager.cleanup_ = func(conn IWSConn) {
			cleanupCalled = true
		}

		wg := sync.WaitGroup{}
		wg.Add(2)

		go func() {
			roomManager.HandleConnection(conn)
			wg.Done()
		}()

		go func() {
			sendCloseMessage(conn)
			wg.Done()
		}()

		wg.Wait()
		assert.True(t, cleanupCalled)
	})

	t.Run("Returns when reading message fails", func(t *testing.T) {
		conn := NewMockChanConn()
		timeout := gotils.NewTimeoutMs(100)
		go func() { panic(<-timeout.ErrorCh) }()
		defer timeout.Cancel()
		roomManager := newRoomManager()
		roomManager.cleanup_ = func(conn IWSConn) {}

		wg := sync.WaitGroup{}
		wg.Add(2)

		go func() {
			roomManager.HandleConnection(conn)
			wg.Done()
		}()

		go func() {
			conn.ReadChan <- MockMessage{
				MessageType: websocket.TextMessage,
				Data:        []byte{},
				Err:         errors.New("FATAL_ERROR"),
			}
			wg.Done()
		}()

		wg.Wait()
	})

	t.Run("Returns ok when handles action succesfully", func(t *testing.T) {
		conn := NewMockChanConn()
		timeout := gotils.NewTimeoutMs(100)
		go func() { panic(<-timeout.ErrorCh) }()
		defer timeout.Cancel()
		manager := newRoomManager()
		manager.actions = actionCollection{
			MessageTypeJoin: func(r *roomManager, conn IWSConn, msg *messageStruct) messageStruct {
				return okMessage
			},
		}
		manager.cleanup_ = func(conn IWSConn) {}

		message := &messageStruct{
			Type: MessageTypeJoin,
		}
		msgData, _ := json.Marshal(message)

		wg := sync.WaitGroup{}
		wg.Add(2)

		go func() {
			manager.HandleConnection(conn)
			wg.Done()
		}()

		go func() {
			conn.ReadChan <- MockMessage{
				MessageType: websocket.TextMessage,
				Data:        msgData,
				Err:         nil,
			}

			receivedMsg := <-conn.WriteChan
			json.Unmarshal(receivedMsg.Data, message)
			assert.Equal(t, MessageTypeOK, message.Type)
			sendCloseMessage(conn)
			wg.Done()
		}()

		wg.Wait()
	})

	t.Run("Returns error when validation fails", func(t *testing.T) {
		conn := NewMockChanConn()
		timeout := gotils.NewTimeoutMs(100)
		go func() { panic(<-timeout.ErrorCh) }()
		defer timeout.Cancel()
		manager := newRoomManager()
		manager.cleanup_ = func(conn IWSConn) {}

		message := &messageStruct{
			Type: "asd",
		}
		msgData, _ := json.Marshal(message)

		wg := sync.WaitGroup{}
		wg.Add(2)

		go func() {
			manager.HandleConnection(conn)
			wg.Done()
		}()

		go func() {
			conn.ReadChan <- MockMessage{
				MessageType: websocket.TextMessage,
				Data:        msgData,
				Err:         nil,
			}

			receivedMsg := <-conn.WriteChan
			json.Unmarshal(receivedMsg.Data, message)
			assert.Equal(t, MessageTypeBadMessage, message.Type)
			sendCloseMessage(conn)
			wg.Done()
		}()

		wg.Wait()
	})
}

func sendCloseMessage(conn *MockChanConn) {
	conn.ReadChan <- MockMessage{
		MessageType: websocket.CloseMessage,
		Data:        []byte{},
		Err:         nil,
	}
}
