package room_manager

import (
	ws_mocks "attiladudas/backend/ws/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCleanup(t *testing.T) {
	t.Run("Deletes room when last player exits", func(t *testing.T) {
		manager := newRoomManager()
		roomName := "room1"
		conn := &ws_mocks.MockChanConn{}
		room := &room{
			name: roomName,
			players: []*player{
				{conn: conn},
			},
		}

		manager.roomsByName[roomName] = room
		manager.roomsByConnection[conn] = room

		manager.cleanup(conn)

		_, ok := manager.roomsByName[roomName]
		assert.False(t, ok)

		_, ok = manager.roomsByConnection[conn]
		assert.False(t, ok)
	})

	t.Run("Removes leaving player from room", func(t *testing.T) {
		manager := newRoomManager()
		roomName := "room1"
		conn1 := &ws_mocks.MockChanConn{}
		conn2 := &ws_mocks.MockChanConn{}
		room := &room{
			name: roomName,
			players: []*player{
				{conn: conn1},
				{conn: conn2},
			},
		}

		manager.roomsByName[roomName] = room
		manager.roomsByConnection[conn1] = room
		manager.roomsByConnection[conn2] = room

		manager.cleanup(conn1)

		_, ok := manager.roomsByName[roomName]
		assert.True(t, ok)

		_, ok = manager.roomsByConnection[conn1]
		assert.False(t, ok)

		_, ok = manager.roomsByConnection[conn2]
		assert.True(t, ok)
	})

	t.Run("Does nothing when no room found", func(t *testing.T) {
		manager := newRoomManager()
		conn := &ws_mocks.MockChanConn{}

		manager.cleanup(conn)
	})
}
