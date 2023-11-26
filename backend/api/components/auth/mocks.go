package auth

type MockJwt struct {
	Encode_ func(*Claims) (string, error)
	Decode_ func(string) (*Claims, error)
}

func (m *MockJwt) Encode(c *Claims) (string, error) {
	return m.Encode_(c)
}

func (m *MockJwt) Decode(s string) (*Claims, error) {
	return m.Decode_(s)
}
