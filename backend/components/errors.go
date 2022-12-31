package components

type ApiError string

func (e ApiError) Error() string {
	return string(e)
}
