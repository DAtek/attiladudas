package fiar

type IGame interface {
	Move(playerName string, pos *Position) error
	GetWinner() *string
	GetSquares() [][]Square
	CurrentPlayer() string
}

type Position struct {
	X uint8
	Y uint8
}

type Square uint8

const (
	SquareEmpty   = Square(0)
	SquarePlayer1 = Square(1)
	SquarePlayer2 = Square(2)
)

type game struct {
	player1       string
	player2       string
	currentPlayer *string
	winner        *string
	squares       [][]Square
}
