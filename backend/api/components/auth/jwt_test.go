package auth

import (
	"crypto"
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func TestJWT(t *testing.T) {
	privateKey := []byte(`-----BEGIN PRIVATE KEY-----
MC4CAQAwBQYDK2VwBCIEIACAVG3FFfs9wr2asYt90wfKLnGOrdEiZziDAAyffR5c
-----END PRIVATE KEY-----
`)

	publicKey := []byte(`-----BEGIN PUBLIC KEY-----
MCowBQYDK2VwAyEAbfTb1EmHuFvX+Ek0tqIpa/45AFKcJYMjY8+VU7hWExo=
-----END PUBLIC KEY-----
`)

	t.Run("Can't create JWTContext with invalid private key", func(t *testing.T) {
		_, err := NewJwtContext([]byte(""), publicKey)

		assert.Error(t, err)
	})

	t.Run("Can't create JWTContext with invalid public key", func(t *testing.T) {
		_, err := NewJwtContext(privateKey, []byte(""))

		assert.Error(t, err)
	})

	t.Run("Encoding and decoding works fine", func(t *testing.T) {
		jwtContext, _ := NewJwtContext(privateKey, publicKey)

		data := &Claims{
			Exp:      1234567890,
			Username: "John Doe",
		}

		token, encodeErr := jwtContext.Encode(data)
		decodedData, decodeErr := jwtContext.Decode(token)

		assert.Nil(t, encodeErr)
		assert.Nil(t, decodeErr)
		assert.Equal(t, data, decodedData)
	})

	t.Run("Decoding returns error when token is invalid", func(t *testing.T) {
		jwtContext, _ := NewJwtContext(privateKey, publicKey)

		_, err := jwtContext.Decode("asd")

		assert.Error(t, err)
	})

	t.Run("Decoding returns error when claims are invalid", func(t *testing.T) {
		jwtContext, _ := NewJwtContext(privateKey, publicKey)
		mockContext := newMockJwtContext(privateKey, publicKey)

		data := jwt.MapClaims{
			"color": "red",
		}

		token, encodeErr := mockContext.Encode(data)
		_, decodeErr := jwtContext.Decode(token)

		assert.Nil(t, encodeErr)
		assert.EqualError(t, InvalidClaimsError, decodeErr.Error())
	})
}

type mockJwtContext struct {
	privateKey crypto.PrivateKey
	publicKey  crypto.PublicKey
}

func newMockJwtContext(pemPrivateKey []byte, pemPublicKey []byte) *mockJwtContext {
	privateKey, _ := jwt.ParseEdPrivateKeyFromPEM(pemPrivateKey)
	publicKey, _ := jwt.ParseEdPublicKeyFromPEM(pemPublicKey)
	return &mockJwtContext{privateKey, publicKey}
}

func (context *mockJwtContext) Encode(claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(&jwt.SigningMethodEd25519{}, claims)
	return token.SignedString(context.privateKey)
}
