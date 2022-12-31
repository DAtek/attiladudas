package app_auth

import (
	"attiladudas/backend/components"
	"attiladudas/backend/components/auth"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

const secondsInMinutes = 60
const tokenExpirationSeconds = 15 * secondsInMinutes

type postTokenBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type tokenResponse struct {
	Token string `json:"token"`
}

func PostTokenHandler(userStore components.IUserStore, jwtContext auth.IJwt) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		postToken(ctx, userStore, jwtContext)
	}
}

func postToken(ctx *gin.Context, userStore components.IUserStore, jwtContext auth.IJwt) {
	request := &postTokenBody{}
	jsonError := ctx.ShouldBindJSON(request)

	if jsonError != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	user, err := userStore.GetUser(request.Username)

	switch err {
	case components.NotFoundError:
		ctx.Status(http.StatusBadRequest)
		return
	case nil:
	default:
		ctx.Status(http.StatusInternalServerError)
		return
	}

	decoded, unexpectedErr := base64.StdEncoding.DecodeString(user.PasswordHash)

	if unexpectedErr != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword(decoded, []byte(request.Password))

	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	token, _ := jwtContext.Encode(&auth.Claims{
		Username: user.Username,
		Exp:      uint(time.Now().Unix() + tokenExpirationSeconds),
	})

	ctx.JSON(http.StatusCreated, tokenResponse{Token: token})
}
