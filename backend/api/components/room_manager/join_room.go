package room_manager

import (
	"encoding/json"

	"github.com/DAtek/golidator"
	"github.com/gorilla/websocket"
)

type joinRoomData struct {
	Player string `json:"player"`
	Room   string `json:"room"`
}

func joinRoom(manager *roomManager, conn IWSConn, msg *messageStruct) messageStruct {
	data, err := parseData(&joinRoomData{}, msg)
	if err != nil {
		return messageStruct{
			Type: MessageTypeBadMessage,
			Data: *err,
		}
	}

	room_, ok := manager.roomsByName[data.Room]
	if !ok {
		players := []*player{
			{
				name: data.Player,
				conn: conn,
			},
		}

		room_ = &room{
			name:    data.Room,
			players: players,
		}

		manager.roomsByName[data.Room] = room_
		manager.roomsByConnection[conn] = room_
		return okMessage
	}

	if err := golidator.Validate(data, room_); err != nil {
		return createBadMessageFromValidationError(err)
	}

	room_.players = append(
		room_.players,
		&player{
			name: data.Player,
			conn: conn,
		},
	)

	manager.roomsByConnection[conn] = room_

	if otherPlayer := room_.getOtherPlayer(data.Player); otherPlayer != nil {
		response, _ := json.Marshal(msg)
		otherPlayer.conn.WriteMessage(websocket.TextMessage, response)
	}
	return okMessage
}

func (obj *joinRoomData) GetValidators(params ...interface{}) golidator.ValidatorCollection {
	room := params[0].(*room)

	return golidator.ValidatorCollection{
		{Field: "room", Function: func() *golidator.ValueError {
			if len(room.players) >= 2 {
				return &golidator.ValueError{
					ErrorType: "ROOM_IS_FULL",
				}
			}
			return nil
		}},
		{Field: "player", Function: func() *golidator.ValueError {
			if len(room.players) == 1 && room.players[0].name == obj.Player {
				return &golidator.ValueError{
					ErrorType: "PLAYER_WITH_THIS_NAME_ALREADY_JOINED",
				}
			}
			return nil
		}},
	}
}
