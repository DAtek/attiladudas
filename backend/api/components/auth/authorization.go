package auth

import (
	"strings"
)

type AuthError string

func (e AuthError) Error() string {
	return string(e)
}

const (
	InvalidAuthHeaderError = AuthError("INVALID OR MISSING AUTHORIZATION HEADER")
	MissingUsernameError   = AuthError("MISSING USERNAME IN JWT")
)

type authContext struct {
	jtwCtx IJwt
}

type IAuthorization interface {
	RequireUsername(authHeader string) error
}

func NewAuthContext(jwtCtx IJwt) IAuthorization {
	return &authContext{
		jtwCtx: jwtCtx,
	}
}

func (a *authContext) RequireUsername(authHeader string) error {
	parts := strings.Split(authHeader, "Bearer ")

	if len(parts) != 2 {
		return InvalidAuthHeaderError
	}

	token := parts[1]
	claims, err := a.jtwCtx.Decode(token)

	if err != nil {
		return err
	}

	if claims.Username == "" {
		return MissingUsernameError
	}

	return nil
}
