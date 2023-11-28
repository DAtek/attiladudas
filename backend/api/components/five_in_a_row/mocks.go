package fiar

type MockGame struct {
	Join_          func(playerName string) error
	Move_          func(playerName string, pos *Position) error
	CurrentPlayer_ func() string
	GetWinner_     func() *string
	GetSquares_    func() [][]Square
}

func (m *MockGame) Join(playerName string) error {
	return m.Join_(playerName)
}

func (m *MockGame) Move(playerName string, pos *Position) error {
	return m.Move_(playerName, pos)
}

func (m *MockGame) GetWinner() *string {
	return m.GetWinner_()
}

func (m *MockGame) GetSquares() [][]Square {
	return m.GetSquares_()
}

func (m *MockGame) CurrentPlayer() string {
	return m.CurrentPlayer_()
}

func newGame() *game {
	g := &game{
		player1: "a",
		player2: "b",
		squares: createSquares(&Size{20, 20}),
	}
	g.currentPlayer = &g.player1
	return g
}
