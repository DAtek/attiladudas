package room_manager

import (
	"attiladudas/backend/ws"
	"encoding/json"

	"github.com/gorilla/websocket"
)

var sendMessage action = func(manager *roomManager, conn ws.IConn, msg *messageStruct) messageStruct {
	room, ok := manager.roomsByConnection[conn]
	if !ok {
		return messageStruct{Type: MessageTypeBadMessage, Data: "PLAYER_DID_NOT_JOIN_ANY_ROOM"}
	}

	var player *player
	for _, player_ := range room.players {
		if player_.conn == conn {
			player = player_
		}
	}

	if otherPlayer := room.getOtherPlayer(player.name); otherPlayer != nil {
		response, _ := json.Marshal(msg)
		otherPlayer.conn.WriteMessage(websocket.TextMessage, response)
	}

	return okMessage
}
