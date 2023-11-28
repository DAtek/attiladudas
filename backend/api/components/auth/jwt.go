package auth

import (
	"crypto"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	Exp      uint   `json:"exp"`
	Username string `json:"username"`
}

type JWTError string

func (e JWTError) Error() string {
	return string(e)
}

const InvalidClaimsError = JWTError("INVALID CLAIMS")

type IJwt interface {
	Encode(*Claims) (string, error)
	Decode(string) (*Claims, error)
}

type jwtContext struct {
	privateKey crypto.PrivateKey
	publicKey  crypto.PublicKey
}

func NewJwtContext(pemPrivateKey []byte, pemPublicKEy []byte) (IJwt, error) {
	privateKey, privateErr := jwt.ParseEdPrivateKeyFromPEM(pemPrivateKey)
	if privateErr != nil {
		return nil, privateErr
	}

	publicKey, publicErr := jwt.ParseEdPublicKeyFromPEM(pemPublicKEy)
	if publicErr != nil {
		return nil, publicErr
	}

	return &jwtContext{privateKey, publicKey}, nil
}

func (context *jwtContext) Encode(claims *Claims) (string, error) {
	token := jwt.NewWithClaims(&jwt.SigningMethodEd25519{}, claims)
	return token.SignedString(context.privateKey)
}

func (context *jwtContext) Decode(encoded string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(encoded, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return context.publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims := token.Claims.(*Claims)
	return claims, nil
}

func (claims Claims) Valid() error {
	if claims.Exp == 0 || claims.Username == "" {
		return InvalidClaimsError
	}

	return nil
}
