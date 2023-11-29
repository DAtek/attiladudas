package room_manager

import (
	"encoding/json"
	"github.com/DAtek/gotils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJoinRoom(t *testing.T) {
	t.Run("Room is being created when player1 joins", func(t *testing.T) {
		manager := newRoomManager()
		playerName := "player1"
		roomName := "room1"
		conn := &MockChanConn{}

		expectedRoom := &room{
			name: roomName,
			players: []*player{
				{
					name: playerName,
					conn: conn,
				},
			},
		}

		data := &joinRoomData{
			Player: playerName,
			Room:   roomName,
		}
		message := &messageStruct{}
		message.setData(data, MessageTypeJoin)

		result := joinRoom(manager, conn, message)

		assert.Equal(t, okMessage, result)

		room := manager.roomsByName[roomName]
		assert.Equal(t, expectedRoom, room)

		room = manager.roomsByConnection[conn]
		assert.Equal(t, expectedRoom, room)
	})

	t.Run("Player1 receives message when player2 joins", func(t *testing.T) {
		timeout := gotils.NewTimeoutMs(100)
		go func() { panic(<-timeout.ErrorCh) }()
		defer timeout.Cancel()
		manager := newRoomManager()
		roomName := "room1"
		conn := NewMockChanConn()
		player1 := "Emily"
		player2 := "Anna"

		manager.roomsByName[roomName] = &room{
			players: []*player{
				{
					conn: conn,
					name: player1,
				},
			},
		}

		request := &joinRoomData{
			Player: player2,
			Room:   roomName,
		}
		message := &messageStruct{}
		message.setData(request, MessageTypeJoin)

		group := gotils.NewGoroGroup()

		group.Add(func() {
			result := joinRoom(manager, nil, message)
			assert.Equal(t, okMessage, result)
		})

		group.Add(func() {
			receivedMessage := <-conn.WriteChan
			msgStruct := &messageStruct{}
			json.Unmarshal(receivedMessage.Data, msgStruct)
			assert.Equal(t, MessageTypeJoin, msgStruct.Type)
			data := &joinRoomData{}
			json.Unmarshal([]byte(msgStruct.Data), data)
			assert.Equal(t, request, data)
		})

		group.Run()
	})

	t.Run("Returns error message when room is full", func(t *testing.T) {
		manager := newRoomManager()
		roomName := "room1"
		manager.roomsByName[roomName] = &room{
			players: []*player{{}, {}},
		}

		message := &messageStruct{}
		message.setData(&joinRoomData{
			Player: "a",
			Room:   roomName,
		}, MessageTypeJoin)

		result := joinRoom(manager, nil, message)

		assert.Equal(t, MessageTypeBadMessage, result.Type)
	})

	t.Run("Returns error message a player already joined with the same name", func(t *testing.T) {
		manager := newRoomManager()
		roomName := "room1"
		playerName := "a"

		manager.roomsByName[roomName] = &room{
			players: []*player{{name: playerName}},
		}

		message := &messageStruct{}
		message.setData(&joinRoomData{
			Player: playerName,
			Room:   roomName,
		}, MessageTypeJoin)

		result := joinRoom(manager, nil, message)

		assert.Equal(t, MessageTypeBadMessage, result.Type)
	})

	t.Run("Returns error if data is invalid", func(t *testing.T) {
		manager := newRoomManager()
		conn1 := NewMockChanConn()

		msg := &messageStruct{
			Type: MessageTypeJoin,
			Data: "",
		}

		result := joinRoom(manager, conn1, msg)

		assert.Equal(t, MessageTypeBadMessage, result.Type)
		assert.Equal(t, "INVALID_DATA", result.Data)
	})
}
