package fiar

type FiveInARowError string

func (e FiveInARowError) Error() string {
	return string(e)
}
