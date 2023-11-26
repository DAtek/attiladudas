package auth

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequireUsername(t *testing.T) {
	auth := &AuthContext{}

	t.Run("Ok if jwt contains username", func(t *testing.T) {
		jwt := &mockJwt{}
		token := "fake_token"
		jwt.decode = func(s string) (*Claims, error) {
			if s == token {
				return &Claims{Username: "Max"}, nil
			}

			return nil, errors.New("WRONG TOKEN")
		}

		err := auth.RequireUsername("Bearer "+token, jwt)

		assert.Nil(t, err)
	})

	t.Run("Error if invalid auth header", func(t *testing.T) {
		jwt := &mockJwt{}

		err := auth.RequireUsername("Bearer<fake_token>", jwt)

		assert.Error(t, err)
	})

	t.Run("Error if invalid jwt", func(t *testing.T) {
		jwt := &mockJwt{}
		jwt.decode = func(s string) (*Claims, error) {
			return nil, errors.New("WRONG TOKEN")
		}

		err := auth.RequireUsername("Bearer asdasda", jwt)

		assert.Error(t, err)
	})

	t.Run("Error if empty username", func(t *testing.T) {
		jwt := &mockJwt{}
		jwt.decode = func(s string) (*Claims, error) {
			return &Claims{Username: ""}, nil
		}

		err := auth.RequireUsername("Bearer asdasda", jwt)

		assert.EqualError(t, MissingUsernameError, err.Error())
	})
}
