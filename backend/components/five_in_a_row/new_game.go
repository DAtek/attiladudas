package fiar

type Size struct {
	X uint8
	Y uint8
}

func NewGame(size *Size, player1, player2 string) IGame {
	game := &game{
		squares: createSquares(size),
		player1: player1,
		player2: player2,
	}

	game.currentPlayer = &game.player1
	return game
}
