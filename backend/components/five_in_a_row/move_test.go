package fiar

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMove(t *testing.T) {
	t.Run("Player1 moves", func(t *testing.T) {
		game := newGame()

		position := &Position{2, 3}
		err := game.Move(game.player1, position)

		assert.Nil(t, err)
		assert.Equal(t, SquarePlayer1, game.squares[position.X][position.Y])
		assert.Equal(t, game.player2, *game.currentPlayer)
		assert.Nil(t, game.winner)
	})

	t.Run("Player2 moves", func(t *testing.T) {
		game := newGame()
		game.currentPlayer = &game.player2

		position := &Position{1, 2}
		err := game.Move(game.player2, position)

		assert.Nil(t, err)
		assert.Equal(t, SquarePlayer2, game.squares[position.X][position.Y])
		assert.Equal(t, game.player1, *game.currentPlayer)
		assert.Nil(t, game.winner)
	})

	t.Run("Moving to an occupied square is forbidden", func(t *testing.T) {
		game := newGame()
		game.currentPlayer = &game.player2
		game.squares[1][2] = 1

		position := &Position{1, 2}
		err := game.Move(game.player2, position)

		assert.EqualError(t, err, ErrorInvalidPosition.Error())
	})

	winnerScenarios := []struct {
		name      string
		positions [][]int
		player    string
	}{
		{
			name:      "up:down",
			positions: [][]int{{10, 9}, {10, 8}, {10, 11}, {10, 12}},
			player:    "a",
		},
		{
			name:      "left:right",
			positions: [][]int{{9, 10}, {8, 10}, {7, 10}, {11, 10}},
			player:    "b",
		},
		{
			name:      "up-right:down-left",
			positions: [][]int{{11, 9}, {12, 8}, {9, 11}, {8, 12}},
			player:    "a",
		},
		{
			name:      "up-left:down-right",
			positions: [][]int{{9, 9}, {8, 8}, {11, 11}, {12, 12}},
			player:    "a",
		},
	}
	for _, scenario := range winnerScenarios {
		t.Run(fmt.Sprintf("Winner is being calculated %s", scenario.name), func(t *testing.T) {
			game := newGame()
			game.currentPlayer = &scenario.player

			for _, pos := range scenario.positions {
				val := game.getSquareValue(scenario.player)
				game.squares[pos[0]][pos[1]] = val
			}

			position := &Position{10, 10}
			game.Move(scenario.player, position)

			assert.Equal(t, scenario.player, *game.winner)
		})
	}

	t.Run("Only current player can move", func(t *testing.T) {
		game := newGame()

		position := &Position{2, 3}
		err := game.Move(game.player2, position)

		assert.EqualError(t, err, ErrorNotYourTurn.Error())
	})

	t.Run("Returns error if X is out of range", func(t *testing.T) {
		game := newGame()
		invalidX := len(game.squares)

		err := game.Move(*game.currentPlayer, &Position{X: uint8(invalidX), Y: 4})

		assert.EqualError(t, err, ErrorInvalidPosition.Error())
	})

	t.Run("Returns error if Y is out of range", func(t *testing.T) {
		game := newGame()
		invalidY := len(game.squares[0])

		err := game.Move(*game.currentPlayer, &Position{X: 0, Y: uint8(invalidY)})

		assert.EqualError(t, err, ErrorInvalidPosition.Error())
	})

	t.Run("Move is not allowed when game ended", func(t *testing.T) {
		game := newGame()
		game.winner = &game.player1

		err := game.Move(*game.currentPlayer, &Position{X: 0, Y: 0})

		assert.EqualError(t, err, ErrorGameEnded.Error())
	})
}
