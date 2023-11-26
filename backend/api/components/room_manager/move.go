package room_manager

import (
	fiar "attiladudas/backend/components/five_in_a_row"
	"attiladudas/backend/ws"

	"github.com/DAtek/golidator"
)

type moveMessageData struct {
	Position []int `json:"position"`
	position *fiar.Position
}

var move action = func(manager *roomManager, conn ws.IConn, msg *messageStruct) messageStruct {
	room, ok := manager.roomsByConnection[conn]
	if !ok || room.game == nil {
		return messageStruct{Type: MessageTypeBadMessage, Data: "NO_ROOM"}
	}

	data, err := parseData(&moveMessageData{}, msg)

	if err != nil {
		return messageStruct{
			Type: MessageTypeBadMessage,
			Data: *err,
		}
	}

	if winner := room.game.GetWinner(); winner != nil {
		return messageStruct{Type: MessageTypeBadMessage, Data: "GAME_ALREADY_ENDED"}
	}

	if err := golidator.Validate(data, room.game); err != nil {
		return createBadMessageFromValidationError(err)
	}

	player := room.getPlayerByConnection(conn)

	if err := room.game.Move(player.name, data.position); err != nil {
		return messageStruct{Type: MessageTypeBadMessage, Data: err.Error()}
	}

	room.sendUpdateGameMessage()
	return okMessage
}

func (obj *moveMessageData) GetValidators(params ...interface{}) golidator.ValidatorCollection {
	game := params[0].(fiar.IGame)
	squares := game.GetSquares()

	return golidator.ValidatorCollection{
		{Field: "position", Function: func() *golidator.ValueError {
			if len(obj.Position) != 2 {
				return &golidator.ValueError{
					ErrorType: "INVALID_POSITION",
					Context:   map[string]any{"required_items": 2},
				}
			}

			x := obj.Position[0]
			y := obj.Position[1]

			if x < 0 || y < 0 || x >= len(squares) || y >= len(squares[0]) {
				return &golidator.ValueError{
					ErrorType: "INVALID_POSITION",
				}
			}

			obj.position = &fiar.Position{
				X: uint8(x),
				Y: uint8(y),
			}

			return nil
		}},
	}
}
