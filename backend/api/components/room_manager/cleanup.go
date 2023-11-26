package room_manager

import (
	"attiladudas/backend/ws"
)

func (manager *roomManager) cleanup(conn ws.IConn) {
	room, ok := manager.roomsByConnection[conn]

	if !ok {
		return
	}

	delete(manager.roomsByConnection, conn)

	remainingPlayers := []*player{}
	for _, player := range room.players {
		if player.conn != conn {
			remainingPlayers = append(remainingPlayers, player)
		}
	}
	room.players = remainingPlayers

	if len(room.players) == 0 {
		delete(manager.roomsByName, room.name)
	}
}
