package fiar

func (g *game) GetSquares() [][]Square {
	squares := [][]Square{}
	for i := 0; i < len(g.squares); i++ {
		col := []Square{}
		for j := 0; j < len(g.squares[0]); j++ {
			col = append(col, g.squares[i][j])
		}
		squares = append(squares, col)
	}

	return squares
}
