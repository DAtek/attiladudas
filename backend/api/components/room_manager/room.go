package room_manager

import (
	fiar "api/components/five_in_a_row"
	"github.com/gofiber/contrib/websocket"
)

type room struct {
	name    string
	players []*player
	game    fiar.IGame
}

type player struct {
	name string
	side playerSide
	conn IWSConn
}

type playerSide string

const (
	sideX = playerSide("X")
	sideO = playerSide("O")
)

var validSides = []playerSide{sideX, sideO}

func (r *room) sendUpdateGameMessage() {
	updateGameMessage := newGameUpdateMessageFromGame(r.game)
	for _, p := range r.players {
		p.conn.WriteMessage(websocket.TextMessage, updateGameMessage)
	}
}

func (r *room) getOtherPlayer(name string) *player {
	if len(r.players) != 2 {
		return nil
	}

	if name == r.players[0].name {
		return r.players[1]
	}

	return r.players[0]
}

func (r *room) getPlayerByConnection(conn IWSConn) *player {
	for _, p := range r.players {
		if conn == p.conn {
			return p
		}
	}

	// this should never happen, joinRoom() is responsible for storing the connections properly
	panic("PLAYER_NOT_EXISTS_IN_ROOM")
}

func (r *room) createNewGame() {
	player1, player2 := r.getPlayersForGame()
	r.game = fiar.NewGame(
		&fiar.Size{X: 11, Y: 11},
		*player1,
		*player2,
	)
}

func (r *room) getPlayersForGame() (*string, *string) {
	if r.players[0].side == sideX {
		return &r.players[0].name, &r.players[1].name
	}
	return &r.players[1].name, &r.players[0].name
}
