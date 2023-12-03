package fiar

func createSquares(size *Size) [][]Square {
	squares := [][]Square{}

	for col := uint8(0); col < size.X; col++ {
		squaresInRow := []Square{}
		for row := uint8(0); row < size.Y; row++ {
			squaresInRow = append(squaresInRow, SquareEmpty)
		}
		squares = append(squares, squaresInRow)
	}
	return squares
}
