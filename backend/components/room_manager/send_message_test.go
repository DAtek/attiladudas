package room_manager

import (
	"attiladudas/backend/helpers"
	ws_mocks "attiladudas/backend/ws/mocks"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendMessage(t *testing.T) {
	t.Run("Other player receives message", func(t *testing.T) {
		timeout := helpers.NewTimeout(100)
		defer timeout.Finish()

		manager := newRoomManager()
		player1 := &player{
			conn: ws_mocks.NewMockChanConn(),
		}

		conn2 := ws_mocks.NewMockChanConn()
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

		wg := helpers.NewWaitGroup()

		wg.Add(func() {
			sendMessage(manager, player1.conn, msg)
		})

		wg.Add(func() {
			receivedMsg := <-conn2.WriteChan
			msgStruct := &messageStruct{}
			json.Unmarshal(receivedMsg.Data, msgStruct)
			assert.Equal(t, msg, msgStruct)
		})

		wg.Wait()
	})

	t.Run("Senging message is forbidden without joining a room", func(t *testing.T) {
		timeout := helpers.NewTimeout(100)
		defer timeout.Finish()
		conn := ws_mocks.NewMockChanConn()
		manager := newRoomManager()

		msg := &messageStruct{
			Type: MessageTypeSendMessage,
			Data: "Hi",
		}

		receivedMsg := sendMessage(manager, conn, msg)

		assert.Equal(t, MessageTypeBadMessage, receivedMsg.Type)

	})
}
