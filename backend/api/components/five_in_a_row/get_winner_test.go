package fiar

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetWinner(t *testing.T) {
	t.Run("Returns winner if there is a winner", func(t *testing.T) {
		game := newGame()
		winner := "player1"
		game.winner = &winner

		actualWinner := game.GetWinner()

		assert.Equal(t, winner, *actualWinner)
	})

	t.Run("Returns nil if nobody won the game", func(t *testing.T) {
		game := newGame()

		winner := game.GetWinner()

		assert.Nil(t, winner)
	})

	t.Run("game.winner is immutable from outside", func(t *testing.T) {
		game := newGame()
		originalWinner := "player1"
		winner := originalWinner
		winner2 := "player2"
		game.winner = &winner

		gameWinner := game.GetWinner()
		*gameWinner = winner2

		assert.Equal(t, originalWinner, *game.winner)
	})
}
