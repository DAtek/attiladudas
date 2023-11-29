package room_manager

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"

	fiar "api/components/five_in_a_row"

	"github.com/DAtek/gotils"
	"github.com/stretchr/testify/assert"
)

func TestMove(t *testing.T) {
	t.Run("IGame move is being called and game update is being sent to both players", func(t *testing.T) {
		timeout := gotils.NewTimeoutMs(100)
		go func() { panic(<-timeout.ErrorCh) }()
		defer timeout.Cancel()

		manager := newRoomManager()
		conn1 := NewMockChanConn()
		conn2 := NewMockChanConn()

		player1 := &player{name: "a", conn: conn1}
		player2 := &player{name: "b", conn: conn2}
		moveCalled := false
		game := &fiar.MockGame{}

		game.Move_ = func(playerName string, pos *fiar.Position) error {
			moveCalled = true
			return nil
		}
		game.GetWinner_ = func() *string {
			return nil
		}
		game.CurrentPlayer_ = func() string { return "" }
		game.GetSquares_ = func() [][]fiar.Square {
			row := []fiar.Square{0, 0, 0}
			return [][]fiar.Square{row, row, row}
		}

		manager.roomsByConnection[conn1] = &room{
			players: []*player{player1, player2},
			game:    game,
		}

		data := &moveMessageData{Position: []int{0, 0}}
		msg := &messageStruct{}
		msg.setData(data, MessageTypeMove)

		wg := gotils.NewGoroGroup()
		wg.Add(func() {
			result := move(manager, conn1, msg)
			assert.Equal(t, okMessage, result)
			assert.True(t, moveCalled)
		})

		wg.Add(func() {
			receivedMessage := <-conn1.WriteChan
			msgStruct := &messageStruct{}
			json.Unmarshal(receivedMessage.Data, msgStruct)
			assert.Equal(t, MessageTypeUpdateGame, msgStruct.Type)

			receivedMessage = <-conn2.WriteChan
			json.Unmarshal(receivedMessage.Data, msgStruct)
			assert.Equal(t, MessageTypeUpdateGame, msgStruct.Type)
		})

		wg.Run()
	})

	t.Run("Returns error if game ended", func(t *testing.T) {
		manager := newRoomManager()
		conn1 := NewMockChanConn()
		conn2 := NewMockChanConn()

		player1 := &player{name: "a", conn: conn1}
		player2 := &player{name: "b", conn: conn2}
		game := &fiar.MockGame{}

		manager.roomsByConnection[conn1] = &room{
			players: []*player{player1, player2},
			game:    game,
		}

		game.GetWinner_ = func() *string {
			winner := ""
			return &winner
		}

		data := &moveMessageData{Position: []int{5, 5}}
		msg := &messageStruct{}
		msg.setData(data, MessageTypeMove)

		result := move(manager, conn1, msg)

		assert.Equal(t, MessageTypeBadMessage, result.Type)
		errorMsg := string(result.Data)
		fmt.Printf("errorMsg: %v\n", errorMsg)
		assert.True(t, strings.Contains(errorMsg, "GAME_ALREADY_ENDED"))
	})

	invalidPositionScenarios := []struct {
		name     string
		position []int
	}{
		{name: "Error if wrong number of positions", position: []int{1}},
		{name: "Error if position has negative value", position: []int{1, -20}},
		{name: "Error if position has too big value", position: []int{3, 2}},
	}
	for _, scenario := range invalidPositionScenarios {
		t.Run(scenario.name, func(t *testing.T) {
			manager := newRoomManager()
			conn1 := NewMockChanConn()
			conn2 := NewMockChanConn()

			player1 := &player{name: "a", conn: conn1}
			player2 := &player{name: "b", conn: conn2}
			game := &fiar.MockGame{}
			game.GetWinner_ = func() *string { return nil }
			game.GetSquares_ = func() [][]fiar.Square {
				row := []fiar.Square{0, 0, 0}
				return [][]fiar.Square{row, row, row}
			}

			manager.roomsByConnection[conn1] = &room{
				players: []*player{player1, player2},
				game:    game,
			}

			data := &moveMessageData{Position: scenario.position}
			msg := &messageStruct{}
			msg.setData(data, MessageTypeMove)

			result := move(manager, conn1, msg)

			assert.Equal(t, MessageTypeBadMessage, result.Type)
			errorMsg := string(result.Data)
			fmt.Printf("errorMsg: %v\n", errorMsg)
			assert.True(t, strings.Contains(errorMsg, "INVALID_POSITION"))
		})
	}

	t.Run("Returns error if no room", func(t *testing.T) {
		manager := newRoomManager()
		conn1 := NewMockChanConn()

		data := &moveMessageData{Position: []int{5, 5}}
		msg := &messageStruct{}
		msg.setData(data, MessageTypeMove)

		result := move(manager, conn1, msg)

		assert.Equal(t, MessageTypeBadMessage, result.Type)
		errorMsg := string(result.Data)
		fmt.Printf("errorMsg: %v\n", errorMsg)
		assert.True(t, strings.Contains(errorMsg, "NO_ROOM"))
	})

	t.Run("Returns error if data is invalid", func(t *testing.T) {
		manager := newRoomManager()
		conn1 := NewMockChanConn()

		player1 := &player{name: "a", conn: conn1}

		manager.roomsByConnection[conn1] = &room{
			players: []*player{player1},
			game:    &fiar.MockGame{},
		}

		msg := &messageStruct{
			Type: MessageTypeMove,
			Data: "asdasd",
		}

		result := move(manager, conn1, msg)

		assert.Equal(t, MessageTypeBadMessage, result.Type)
	})

	t.Run("Returns error core game return error", func(t *testing.T) {
		manager := newRoomManager()
		conn1 := NewMockChanConn()

		player1 := &player{name: "a", conn: conn1}

		errorMessage := "UNEXPECTED_ERROR"
		game := &fiar.MockGame{}
		game.Move_ = func(playerName string, pos *fiar.Position) error {
			return errors.New(errorMessage)
		}
		game.GetWinner_ = func() *string { return nil }
		game.GetSquares_ = func() [][]fiar.Square {
			row := []fiar.Square{0, 0, 0}
			return [][]fiar.Square{row, row, row}
		}

		manager.roomsByConnection[conn1] = &room{
			players: []*player{player1},
			game:    game,
		}

		data := &moveMessageData{Position: []int{0, 0}}
		msg := &messageStruct{}
		msg.setData(data, MessageTypeMove)

		result := move(manager, conn1, msg)

		assert.Equal(t, MessageTypeBadMessage, result.Type)
		assert.Equal(t, errorMessage, result.Data)
	})
}
