package room_manager

import (
	"attiladudas/backend/helpers"
	ws_mocks "attiladudas/backend/ws/mocks"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	fiar "attiladudas/backend/components/five_in_a_row"
	fiar_mocks "attiladudas/backend/components/five_in_a_row/mocks"

	"github.com/stretchr/testify/assert"
)

func TestMove(t *testing.T) {
	t.Run("IGame move is being called and game update is being sent to both players", func(t *testing.T) {
		timeout := helpers.NewTimeout(100)
		defer timeout.Finish()

		manager := newRoomManager()
		conn1 := ws_mocks.NewMockChanConn()
		conn2 := ws_mocks.NewMockChanConn()

		player1 := &player{name: "a", conn: conn1}
		player2 := &player{name: "b", conn: conn2}
		moveCalled := false
		game := &fiar_mocks.MockGame{}

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

		wg := helpers.NewWaitGroup()
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

		wg.Wait()
	})

	t.Run("Returns error if game ended", func(t *testing.T) {
		manager := newRoomManager()
		conn1 := ws_mocks.NewMockChanConn()
		conn2 := ws_mocks.NewMockChanConn()

		player1 := &player{name: "a", conn: conn1}
		player2 := &player{name: "b", conn: conn2}
		game := &fiar_mocks.MockGame{}

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
			conn1 := ws_mocks.NewMockChanConn()
			conn2 := ws_mocks.NewMockChanConn()

			player1 := &player{name: "a", conn: conn1}
			player2 := &player{name: "b", conn: conn2}
			game := &fiar_mocks.MockGame{}
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
}
