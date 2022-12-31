package fiar

import (
	"math"
)

var countMatchFunctions = [][]func(*game, *Position, string) uint8{
	{
		countUp,
		countDown,
	},
	{
		countLeft,
		countRight,
	},
	{
		countUpRight,
		countDownLeft,
	},
	{
		countUpLeft,
		countDownRight,
	},
}

func (g *game) checkWinner(pos *Position, player string) bool {
	for _, countPair := range countMatchFunctions {
		matches := uint8(0)
		for _, count := range countPair {
			matches += count(g, pos, player)
			if matches == 4 {
				return true
			}
		}
	}
	return false
}

func countUp(g *game, pos *Position, player string) uint8 {
	squaresToUp := uint8(math.Min(
		float64(pos.Y),
		4,
	))

	matches := uint8(0)
	for j := uint8(1); j <= squaresToUp; j++ {
		y := pos.Y - j
		if g.squares[pos.X][y] != g.getSquareValue(player) {
			break
		}
		matches++
	}
	return matches
}

func countUpRight(g *game, pos *Position, player string) uint8 {
	squaresToUp := math.Min(
		float64(pos.Y),
		4,
	)

	squaresToUpRight := uint8(math.Min(
		squaresToUp,
		float64(len(g.squares[0]))-float64(pos.X)-1,
	))

	matches := uint8(0)
	for ij := uint8(1); ij <= squaresToUpRight; ij++ {
		y := pos.Y - ij
		x := pos.X + ij
		if g.squares[x][y] != g.getSquareValue(player) {
			break
		}
		matches++
	}
	return matches
}

func countDownLeft(g *game, pos *Position, player string) uint8 {
	squaresToDown := math.Min(
		float64(len(g.squares))-float64(pos.Y)-1,
		4,
	)

	squaresToDownLeft := uint8(math.Min(
		squaresToDown,
		float64(pos.X),
	))

	matches := uint8(0)
	for ij := uint8(1); ij <= squaresToDownLeft; ij++ {
		y := pos.Y + ij
		x := pos.X - ij
		if g.squares[x][y] != g.getSquareValue(player) {
			break
		}
		matches++
	}
	return matches
}

func countUpLeft(g *game, pos *Position, player string) uint8 {
	squaresToUp := math.Min(
		float64(pos.Y),
		4,
	)

	squaresToUpLeft := uint8(math.Min(
		squaresToUp,
		float64(pos.X),
	))

	matches := uint8(0)
	for ij := uint8(1); ij <= squaresToUpLeft; ij++ {
		y := pos.Y - ij
		x := pos.X - ij
		if g.squares[x][y] != g.getSquareValue(player) {
			break
		}
		matches++
	}
	return matches
}

func countDownRight(g *game, pos *Position, player string) uint8 {
	squaresToDown := math.Min(
		float64(len(g.squares))-float64(pos.Y)-1,
		4,
	)

	squaresToDownRight := uint8(math.Min(
		squaresToDown,
		float64(len(g.squares))-float64(pos.X)-1,
	))

	matches := uint8(0)
	for ij := uint8(1); ij <= squaresToDownRight; ij++ {
		y := pos.Y + ij
		x := pos.X + ij
		if g.squares[x][y] != g.getSquareValue(player) {
			break
		}
		matches++
	}
	return matches
}

func countRight(g *game, pos *Position, player string) uint8 {
	squaresToRight := uint8(math.Min(
		float64(len(g.squares[0]))-float64(pos.X)-1,
		4,
	))

	matches := uint8(0)
	for i := uint8(1); i <= squaresToRight; i++ {
		x := pos.X + i
		if g.squares[x][pos.Y] != g.getSquareValue(player) {
			break
		}
		matches++
	}
	return matches
}

func countLeft(g *game, pos *Position, player string) uint8 {
	squaresToLeft := uint8(math.Min(
		float64(pos.X),
		4,
	))

	matches := uint8(0)
	for i := uint8(1); i <= squaresToLeft; i++ {
		x := pos.X - i
		if g.squares[x][pos.Y] != g.getSquareValue(player) {
			break
		}
		matches++
	}
	return matches
}

func countDown(g *game, pos *Position, player string) uint8 {
	squaresToDown := uint8(math.Min(
		float64(len(g.squares))-float64(pos.Y)-1,
		4,
	))
	matches := uint8(0)
	for j := uint8(1); j <= squaresToDown; j++ {
		y := pos.Y + j
		if g.squares[pos.X][y] != g.getSquareValue(player) {
			break
		}
		matches++
	}
	return matches
}
