package token_post

import (
	"api/components/auth"
	"fibertools"
	"time"

	"github.com/DAtek/gotils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const secondsInMinutes = 60
const tokenExpirationSeconds = 15 * secondsInMinutes

type tokenResponse struct {
	Token string `json:"token"`
}

const path = "/api/token/"

func PluginTokenPost(session *gorm.DB, jwtContext auth.IJwt) fibertools.Plugin {
	handler := func(ctx *fiber.Ctx) error {
		return postToken(ctx, session, jwtContext)
	}

	return func(app *fiber.App) {
		app.Post(path, handler)
	}
}

func postToken(ctx *fiber.Ctx, session *gorm.DB, jwtContext auth.IJwt) error {
	data, err := fibertools.BindAndValidateObj[postTokenBody](ctx.BodyParser, session)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fibertools.JsonErrorFromValidationError(err))
	}

	token := gotils.ResultOrPanic(jwtContext.Encode(&auth.Claims{
		Username: data.user.Username,
		Exp:      uint(time.Now().Unix() + tokenExpirationSeconds),
	}))

	ctx.Status(fiber.StatusCreated).JSON(&tokenResponse{Token: token})
	return nil
}
