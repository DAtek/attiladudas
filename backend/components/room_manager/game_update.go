package room_manager

import (
	fiar "attiladudas/backend/components/five_in_a_row"
	"encoding/json"
	"fmt"
	"strings"
)

type gameUpdate struct {
	CurrentPlayer string
	Winner        string
	Squares       [][]fiar.Square
}

// otherwise squares will be serialised as characters
func (obj *gameUpdate) MarshalJSON() ([]byte, error) {
	squaresStringParts := []string{}
	for _, row := range obj.Squares {
		items := []string{}
		for _, square := range row {
			items = append(items, fmt.Sprintf("%d", square))
		}
		itemsString := "[" + strings.Join(items, ",") + "]"
		squaresStringParts = append(squaresStringParts, itemsString)
	}

	squaresString := strings.Join(squaresStringParts, ",")
	squaresString = "[" + squaresString + "]"

	result := fmt.Sprintf(
		`{"currentPlayer":"%s","winner":"%s","squares":%s}`,
		obj.CurrentPlayer,
		obj.Winner,
		squaresString,
	)

	return []byte(result), nil
}

func newGameUpdateMessageFromGame(game fiar.IGame) []byte {
	winner := ""
	if winner_ := game.GetWinner(); winner_ != nil {
		winner = *winner_
	}

	data := &gameUpdate{
		CurrentPlayer: game.CurrentPlayer(),
		Squares:       game.GetSquares(),
		Winner:        winner,
	}
	serializedData, _ := json.Marshal(data)
	msgStuct := &messageStruct{
		Type: MessageTypeUpdateGame,
		Data: string(serializedData),
	}

	msg, _ := json.Marshal(msgStuct)
	return msg
}
