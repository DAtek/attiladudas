package room_manager

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetOtherPlayer(t *testing.T) {
	t.Run("Returns nil if there are less than 2 player", func(t *testing.T) {
		r := room{players: []*player{{}}}

		otherPlayer := r.getOtherPlayer("")

		assert.Nil(t, otherPlayer)
	})

	okScenarios := []struct {
		name           string
		expectedPlayer string
	}{
		{"player1", "player2"},
		{"player2", "player1"},
	}

	for _, scenarion := range okScenarios {
		t.Run(fmt.Sprintf("Returns %s", scenarion.expectedPlayer), func(t *testing.T) {
			r := room{players: []*player{{name: "player1"}, {name: "player2"}}}

			otherPlayer := r.getOtherPlayer(scenarion.name)

			assert.Equal(t, scenarion.expectedPlayer, otherPlayer.name)
		})
	}
}

func TestGetPlayersForGame(t *testing.T) {

	t.Run("Returns player names in correct order 1", func(t *testing.T) {
		r := room{players: []*player{{name: "a", side: sideX}, {name: "b", side: sideO}}}

		player1, player2 := r.getPlayersForGame()

		assert.Equal(t, "a", *player1)
		assert.Equal(t, "b", *player2)
	})

	t.Run("Returns player names in correct order 2", func(t *testing.T) {
		r := room{players: []*player{{name: "a", side: sideO}, {name: "b", side: sideX}}}

		player1, player2 := r.getPlayersForGame()

		assert.Equal(t, "b", *player1)
		assert.Equal(t, "a", *player2)
	})
}
