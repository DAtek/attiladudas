package fiar

func (g *game) GetWinner() *string {
	if g.winner == nil {
		return nil
	}
	winner := *g.winner
	return &winner
}
