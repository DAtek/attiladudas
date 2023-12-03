package auth

type MockAuthContext struct {
	RequireUsername_ func(authHeader string) error
}

func (a *MockAuthContext) RequireUsername(authHeader string) error {
	return a.RequireUsername_(authHeader)
}

type MockJwtContext struct {
	Decode_ func(string) (*Claims, error)
	Encode_ func(*Claims) (string, error)
}

func (c *MockJwtContext) Decode(s string) (*Claims, error) {
	return c.Decode_(s)
}

func (c *MockJwtContext) Encode(claims *Claims) (string, error) {
	return c.Encode_(claims)
}
