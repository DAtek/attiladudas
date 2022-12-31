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
}
