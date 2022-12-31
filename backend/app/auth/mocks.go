package app_auth

import (
	"attiladudas/backend/app/engine"
	"attiladudas/backend/components/auth"
	auth_mocks "attiladudas/backend/components/auth/mocks"
	component_mocks "attiladudas/backend/components/mocks"

	"github.com/gin-gonic/gin"
)

type mockApp struct {
	userStore  *component_mocks.MockUserStore
	jwtContext *auth_mocks.MockJwtContext
	engine     *gin.Engine
}

func newMockApp() *mockApp {
	userStore := &component_mocks.MockUserStore{}

	jwtContext := &auth_mocks.MockJwtContext{
		Encode_: func(c *auth.Claims) (string, error) { return "", nil },
	}

	appEngine := engine.NewEngine(
		nil,
		jwtContext,
		&engine.HandlerCollection{
			PostTokenHandler: PostTokenHandler(userStore, jwtContext),
		},
		"",
	)

	return &mockApp{
		userStore:  userStore,
		jwtContext: jwtContext,
		engine:     appEngine,
	}
}
