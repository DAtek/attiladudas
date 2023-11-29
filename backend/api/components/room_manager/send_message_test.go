package room_manager

import (
	"encoding/json"
	"testing"

	"github.com/DAtek/gotils"
	"github.com/stretchr/testify/assert"
)

func TestSendMessage(t *testing.T) {
	t.Run("Other player receives message", func(t *testing.T) {
		timeout := gotils.NewTimeoutMs(100)
		go func() { panic(<-timeout.ErrorCh) }()
		defer timeout.Cancel()

		manager := newRoomManager()
		player1 := &player{
			conn: NewMockChanConn(),
		}

		conn2 := NewMockChanConn()
		player2 := &player{
			conn: conn2,
		}

		manager.roomsByConnection[player1.conn] = &room{
			players: []*player{
				player1,
				player2,
			},
		}

		msg := &messageStruct{
			Type: MessageTypeSendMessage,
			Data: "Hi",
		}

		wg := gotils.NewGoroGroup()

		wg.Add(func() {
			sendMessage(manager, player1.conn, msg)
		})

		wg.Add(func() {
			receivedMsg := <-conn2.WriteChan
			msgStruct := &messageStruct{}
			json.Unmarshal(receivedMsg.Data, msgStruct)
			assert.Equal(t, msg, msgStruct)
		})

		wg.Run()
	})

	t.Run("Senging message is forbidden without joining a room", func(t *testing.T) {
		timeout := gotils.NewTimeoutMs(100)
		go func() { panic(<-timeout.ErrorCh) }()
		defer timeout.Cancel()
		conn := NewMockChanConn()
		manager := newRoomManager()

		msg := &messageStruct{
			Type: MessageTypeSendMessage,
			Data: "Hi",
		}

		receivedMsg := sendMessage(manager, conn, msg)

		assert.Equal(t, MessageTypeBadMessage, receivedMsg.Type)

	})
}
