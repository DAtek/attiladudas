package auth

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequireUsername(t *testing.T) {

	t.Run("Ok if jwt contains username", func(t *testing.T) {
		jwt := &MockJwtContext{}
		auth := NewAuthContext(jwt)
		token := "fake_token"
		jwt.Decode_ = func(s string) (*Claims, error) {
			if s == token {
				return &Claims{Username: "Max"}, nil
			}

			return nil, errors.New("WRONG TOKEN")
		}

		err := auth.RequireUsername("Bearer " + token)

		assert.Nil(t, err)
	})

	t.Run("Error if invalid auth header", func(t *testing.T) {
		jwt := &MockJwtContext{}
		auth := NewAuthContext(jwt)

		err := auth.RequireUsername("Bearer<fake_token>")

		assert.Error(t, err)
	})

	t.Run("Error if invalid jwt", func(t *testing.T) {
		jwt := &MockJwtContext{}
		auth := NewAuthContext(jwt)
		jwt.Decode_ = func(s string) (*Claims, error) {
			return nil, errors.New("WRONG TOKEN")
		}

		err := auth.RequireUsername("Bearer asdasda")

		assert.Error(t, err)
	})

	t.Run("Error if empty username", func(t *testing.T) {
		jwt := &MockJwtContext{}
		auth := NewAuthContext(jwt)
		jwt.Decode_ = func(s string) (*Claims, error) {
			return &Claims{Username: ""}, nil
		}

		err := auth.RequireUsername("Bearer asdasda")

		assert.EqualError(t, MissingUsernameError, err.Error())
	})
}
