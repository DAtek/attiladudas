package auth_mocks

import (
	"attiladudas/backend/components/auth"
)

type MockAuthContext struct {
	RequireUsername_ func(authHeader string, jwtCtx auth.IJwt) error
}

func (a *MockAuthContext) RequireUsername(authHeader string, j auth.IJwt) error {
	return a.RequireUsername_(authHeader, j)
}
