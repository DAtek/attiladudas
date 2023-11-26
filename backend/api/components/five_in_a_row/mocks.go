package fiar

func newGame() *game {
	g := &game{
		player1: "a",
		player2: "b",
		squares: createSquares(&Size{20, 20}),
	}
	g.currentPlayer = &g.player1
	return g
}
