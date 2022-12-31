package fiar

const ErrorNotYourTurn = FiveInARowError("NOT_YOUR_TURN")
const ErrorInvalidPosition = FiveInARowError("INVALID_POSITION")
const ErrorGameEnded = FiveInARowError("GAME_ENDED")

func (g *game) Move(playerName string, pos *Position) error {
	if g.winner != nil {
		return ErrorGameEnded
	}

	if int(pos.X) >= len(g.squares) || int(pos.Y) >= len(g.squares[0]) {
		return ErrorInvalidPosition
	}

	if g.squares[pos.X][pos.Y] != 0 {
		return ErrorInvalidPosition
	}

	if playerName != *g.currentPlayer {
		return ErrorNotYourTurn
	}

	value := g.getSquareValue(playerName)
	g.squares[pos.X][pos.Y] = value

	if g.checkWinner(pos, playerName) {
		g.winner = &playerName
	}

	if g.currentPlayer == &g.player1 {
		g.currentPlayer = &g.player2
	} else {
		g.currentPlayer = &g.player1
	}

	return nil
}

func (g *game) getSquareValue(playerName string) Square {
	if playerName == g.player1 {
		return SquarePlayer1
	}

	return SquarePlayer2
}
