package auth

type mockJwt struct {
	encode func(*Claims) (string, error)
	decode func(string) (*Claims, error)
}

func (m *mockJwt) Encode(c *Claims) (string, error) {
	return m.encode(c)
}

func (m *mockJwt) Decode(s string) (*Claims, error) {
	return m.decode(s)
}
