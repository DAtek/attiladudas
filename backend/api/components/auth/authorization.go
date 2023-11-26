package auth

import (
	"api/components"
	"strings"
)

const InvalidAuthHeaderError = components.ApiError("INVALID OR MISSING AUTHORIZATION HEADER")
const MissingUsernameError = components.ApiError("MISSING USERNAME IN JWT")

type AuthContext struct{}

type IAuthorization interface {
	RequireUsername(authHeader string, jwtCtx IJwt) error
}

func (a *AuthContext) RequireUsername(authHeader string, jwtCtx IJwt) error {
	parts := strings.Split(authHeader, "Bearer ")

	if len(parts) != 2 {
		return InvalidAuthHeaderError
	}

	token := parts[1]
	claims, err := jwtCtx.Decode(token)

	if err != nil {
		return err
	}

	if claims.Username == "" {
		return MissingUsernameError
	}

	return nil
}
