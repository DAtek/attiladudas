package fiar_mocks

import fiar "attiladudas/backend/components/five_in_a_row"

type MockGame struct {
	Join_          func(playerName string) error
	Move_          func(playerName string, pos *fiar.Position) error
	CurrentPlayer_ func() string
	GetWinner_     func() *string
	GetSquares_    func() [][]fiar.Square
}

func (m *MockGame) Join(playerName string) error {
	return m.Join_(playerName)
}

func (m *MockGame) Move(playerName string, pos *fiar.Position) error {
	return m.Move_(playerName, pos)
}

func (m *MockGame) GetWinner() *string {
	return m.GetWinner_()
}

func (m *MockGame) GetSquares() [][]fiar.Square {
	return m.GetSquares_()
}

func (m *MockGame) CurrentPlayer() string {
	return m.CurrentPlayer_()
}
