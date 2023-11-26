package room_manager

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRoomManager(t *testing.T) {
	t.Run("actions is not nil", func(t *testing.T) {
		manager := NewRoomManager()
		roomManager, _ := manager.(*roomManager)

		assert.NotNil(t, roomManager.actions)
		assert.NotNil(t, roomManager.cleanup_)
	})
}
