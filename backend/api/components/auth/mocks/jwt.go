package auth_mocks

import "attiladudas/backend/components/auth"

type MockJwtContext struct {
	Decode_ func(string) (*auth.Claims, error)
	Encode_ func(*auth.Claims) (string, error)
}

func (c *MockJwtContext) Decode(s string) (*auth.Claims, error) {
	return c.Decode_(s)
}

func (c *MockJwtContext) Encode(claims *auth.Claims) (string, error) {
	return c.Encode_(claims)
}
