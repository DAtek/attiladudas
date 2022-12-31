package fiar

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGame(t *testing.T) {
	t.Run("Creates game with squares", func(t *testing.T) {
		size := &Size{3, 4}

		igame := NewGame(size, "a", "b")
		game, ok := igame.(*game)

		assert.True(t, ok)
		assert.Equal(t, int(size.X), len(game.squares))

		for _, row := range game.squares {
			assert.Equal(t, int(size.Y), len(row))
			for _, square := range row {
				assert.Equal(t, SquareEmpty, square)
			}
		}
	})

	t.Run("Players are properly set", func(t *testing.T) {
		size := &Size{3, 4}
		player1 := "player1"
		player2 := "player2"
		igame := NewGame(size, player1, player2)
		game, ok := igame.(*game)

		assert.True(t, ok)
		assert.Equal(t, game.player1, *game.currentPlayer)
		assert.Equal(t, player1, game.player1)
		assert.Equal(t, player2, game.player2)
	})

	t.Run("Winner is nil", func(t *testing.T) {
		igame := NewGame(&Size{3, 4}, "a", "b")
		game, ok := igame.(*game)

		assert.True(t, ok)
		assert.Nil(t, game.winner)
	})
}
