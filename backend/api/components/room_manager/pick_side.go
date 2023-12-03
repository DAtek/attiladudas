package room_manager

import (
	"encoding/json"

	"github.com/DAtek/golidator"
	"github.com/gorilla/websocket"
)

type pickSideData struct {
	Side playerSide `json:"side"`
}

func pickSide(manager *roomManager, conn IWSConn, msg *messageStruct) messageStruct {
	room, ok := manager.roomsByConnection[conn]

	if !ok {
		return messageStruct{Type: MessageTypeBadMessage, Data: "NO_ROOM"}
	}

	data, err := parseData(&pickSideData{}, msg)
	if err != nil {
		return messageStruct{
			Type: MessageTypeBadMessage,
			Data: *err,
		}
	}

	if err := golidator.Validate(data, room, conn); err != nil {
		return createBadMessageFromValidationError(err)
	}

	player := room.getPlayerByConnection(conn)
	player.side = data.Side

	if room.game == nil && len(room.players) == 2 && room.players[0].side != "" && room.players[1].side != "" {
		room.createNewGame()
	}

	if len(room.players) == 2 && room.players[0].side != "" && room.players[1].side != "" {
		room.sendUpdateGameMessage()
	}

	if otherPlayer := room.getOtherPlayer(player.name); otherPlayer != nil {
		response, _ := json.Marshal(msg)
		otherPlayer.conn.WriteMessage(websocket.TextMessage, response)
	}

	return okMessage
}

func (obj *pickSideData) GetValidators(params ...interface{}) golidator.ValidatorCollection {
	room := params[0].(*room)
	conn := params[1].(IWSConn)

	return golidator.ValidatorCollection{
		{Field: "side", Function: func() *golidator.ValueError {
			if !isSideValid(obj.Side) {
				return &golidator.ValueError{ErrorType: "INVALID_SIDE"}
			}

			if len(room.players) != 2 {
				return &golidator.ValueError{ErrorType: "BOTH_PLAYERS_MUST_JOIN"}
			}
			p := room.getPlayerByConnection(conn)

			if otherPlayer := room.getOtherPlayer(p.name); otherPlayer != nil && otherPlayer.side == obj.Side {
				return &golidator.ValueError{ErrorType: "SIDE_ALREADY_TAKEN"}
			}
			return nil
		}},
	}
}

func isSideValid(side playerSide) bool {
	for _, item := range validSides {
		if side == item {
			return true
		}
	}
	return false
}
