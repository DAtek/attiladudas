package room_manager

import (
	ws_mocks "attiladudas/backend/ws/mocks"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/DAtek/gotils"
	"github.com/stretchr/testify/assert"
)

func TestPickSide(t *testing.T) {
	t.Run("Choosen side is stored properly", func(t *testing.T) {
		timeout := gotils.NewTimeoutMs(100)
		go func() { panic(<-timeout.ErrorCh) }()
		defer timeout.Cancel()

		manager := newRoomManager()
		playerName := "player1"
		p := &player{name: playerName, conn: ws_mocks.NewMockChanConn()}
		conn2 := ws_mocks.NewMockChanConn()

		manager.roomsByConnection[p.conn] = &room{
			players: []*player{
				p,
				{conn: conn2},
			},
		}

		data := &pickSideData{
			Side: sideO,
		}
		message := &messageStruct{}
		message.setData(data, MessageTypePickSide)
		wg := gotils.NewGoroGroup()

		wg.Add(func() {
			result := pickSide(manager, p.conn, message)
			assert.Equal(t, okMessage, result)
			assert.Equal(t, sideO, p.side)
		})

		wg.Add(func() {
			<-conn2.WriteChan
		})

		wg.Run()

	})

	t.Run("Other player gets notified", func(t *testing.T) {
		timeout := gotils.NewTimeoutMs(100)
		go func() { panic(<-timeout.ErrorCh) }()
		defer timeout.Cancel()

		manager := newRoomManager()
		playerName := "player1"
		conn1 := ws_mocks.NewMockChanConn()
		conn2 := ws_mocks.NewMockChanConn()

		manager.roomsByConnection[conn1] = &room{
			players: []*player{
				{name: playerName, conn: conn1},
				{name: "player2", conn: conn2},
			},
		}

		request := &pickSideData{
			Side: sideO,
		}

		message := &messageStruct{}
		message.setData(request, MessageTypePickSide)

		wg := gotils.NewGoroGroup()

		wg.Add(func() {
			result := pickSide(manager, conn1, message)
			assert.Equal(t, okMessage, result)
		})

		wg.Add(func() {
			receivedMessage := <-conn2.WriteChan
			msgStruct := &messageStruct{}
			json.Unmarshal(receivedMessage.Data, msgStruct)
			assert.Equal(t, MessageTypePickSide, msgStruct.Type)
			data := &pickSideData{}
			json.Unmarshal([]byte(msgStruct.Data), data)
			assert.Equal(t, request, data)
		})

		wg.Run()
	})

	t.Run("Game is being created when both players picked side", func(t *testing.T) {
		timeout := gotils.NewTimeoutMs(100)
		defer timeout.Cancel()

		manager := newRoomManager()
		playerName := "player1"
		conn1 := ws_mocks.NewMockChanConn()
		conn2 := ws_mocks.NewMockChanConn()

		room := &room{
			players: []*player{
				{name: playerName, conn: conn1},
				{name: "player2", conn: conn2, side: sideO},
			},
		}
		manager.roomsByConnection[conn1] = room

		request := &pickSideData{
			Side: sideX,
		}

		message := &messageStruct{}
		message.setData(request, MessageTypePickSide)

		wg := gotils.NewGoroGroup()

		wg.Add(func() {
			result := pickSide(manager, conn1, message)
			assert.Equal(t, okMessage, result)
		})

		wg.Add(func() {
			receivedMessage := <-conn1.WriteChan
			msgStruct := &messageStruct{}
			json.Unmarshal(receivedMessage.Data, msgStruct)
			assert.Equal(t, MessageTypeUpdateGame, msgStruct.Type)

			receivedMessage = <-conn2.WriteChan
			json.Unmarshal(receivedMessage.Data, msgStruct)
			assert.Equal(t, MessageTypeUpdateGame, msgStruct.Type)

			<-conn2.WriteChan
		})

		wg.Run()

		assert.Equal(t, playerName, room.game.CurrentPlayer())
	})

	t.Run("Picking the same side is forbidden", func(t *testing.T) {
		manager := newRoomManager()
		playerName := "b"
		conn := ws_mocks.NewMockChanConn()
		manager.roomsByConnection[conn] = &room{
			players: []*player{
				{
					name: "a",
					side: sideX,
				},
				{
					name: playerName,
					conn: conn,
				},
			},
		}

		data := &pickSideData{
			Side: sideX,
		}
		message := &messageStruct{}
		message.setData(data, MessageTypePickSide)

		result := pickSide(manager, conn, message)

		assert.Equal(t, MessageTypeBadMessage, result.Type)
		resultStr := string(result.Data)
		fmt.Printf("resultStr: %v\n", resultStr)
		assert.True(t, strings.Contains(resultStr, "SIDE_ALREADY_TAKEN"))
	})

	t.Run("Picking side is forbidden until both players joined", func(t *testing.T) {
		manager := newRoomManager()
		conn := ws_mocks.NewMockChanConn()

		manager.roomsByConnection[conn] = &room{
			players: []*player{{}},
		}

		data := &pickSideData{
			Side: sideX,
		}
		message := &messageStruct{}
		message.setData(data, MessageTypePickSide)

		result := pickSide(manager, conn, message)

		assert.Equal(t, MessageTypeBadMessage, result.Type)
		resultStr := string(result.Data)
		fmt.Printf("resultStr: %v\n", resultStr)
		assert.True(t, strings.Contains(resultStr, "BOTH_PLAYERS_MUST_JOIN"))
	})

	t.Run("Returns error if side is invalid", func(t *testing.T) {
		manager := newRoomManager()
		conn := ws_mocks.NewMockChanConn()
		manager.roomsByConnection[conn] = &room{
			players: []*player{
				{
					name: "a",
					conn: conn,
				},
				{
					name: "c",
				},
			},
		}

		data := &pickSideData{
			Side: "asd",
		}
		message := &messageStruct{}
		message.setData(data, MessageTypePickSide)

		result := pickSide(manager, conn, message)

		assert.Equal(t, MessageTypeBadMessage, result.Type)
		resultStr := string(result.Data)
		fmt.Printf("resultStr: %v\n", resultStr)
		assert.True(t, strings.Contains(resultStr, "INVALID_SIDE"))
	})

	t.Run("Returns error if message data is invalid", func(t *testing.T) {
		manager := newRoomManager()
		conn := ws_mocks.NewMockChanConn()
		manager.roomsByConnection[conn] = &room{
			players: []*player{
				{
					name: "a",
					conn: conn,
				},
				{
					name: "c",
				},
			},
		}

		message := &messageStruct{Type: MessageTypePickSide, Data: ""}
		result := pickSide(manager, conn, message)

		assert.Equal(t, MessageTypeBadMessage, result.Type)
		resultStr := string(result.Data)
		fmt.Printf("resultStr: %v\n", resultStr)
		assert.True(t, strings.Contains(resultStr, "INVALID_DATA"))
	})

	t.Run("Returns error if no room found", func(t *testing.T) {
		manager := newRoomManager()
		conn := ws_mocks.NewMockChanConn()

		result := pickSide(manager, conn, &messageStruct{})

		assert.Equal(t, MessageTypeBadMessage, result.Type)
		resultStr := string(result.Data)
		fmt.Printf("resultStr: %v\n", resultStr)
		assert.True(t, strings.Contains(resultStr, "NO_ROOM"))
	})

}
