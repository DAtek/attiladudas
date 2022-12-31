package fiar

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSquares(t *testing.T) {
	t.Run("Returns squares", func(t *testing.T) {
		game := newGame()

		squares := game.GetSquares()

		assert.Equal(t, game.squares, squares)
	})

	t.Run("game.squares is immutable from outside", func(t *testing.T) {
		game := newGame()
		size := &Size{5, 5}
		game.squares = createSquares(size)

		squares := game.GetSquares()
		squares[0][0] = SquarePlayer1

		originalSquares := createSquares(size)
		assert.Equal(t, originalSquares, game.squares)
	})
}
