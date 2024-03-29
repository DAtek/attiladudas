package token_post

import (
	"api"
	"db"
	"db/models"
	"encoding/base64"

	"github.com/DAtek/golidator"
	"github.com/DAtek/gotils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type postTokenBody struct {
	Username string `json:"username"`
	Password string `json:"password"`

	user *models.User
}

func (p *postTokenBody) GetValidators(params ...any) golidator.ValidatorCollection {
	session := params[0].(*gorm.DB)

	return golidator.ValidatorCollection{
		{Field: "username", Function: func() *golidator.ValueError {
			if p.Username == "" {
				return api.ErrorRequired
			}
			return nil
		}},
		{Field: "password", Function: func() *golidator.ValueError {
			if p.Password == "" {
				return api.ErrorRequired
			}
			return nil
		}},
		{Field: "__root__", Function: func() *golidator.ValueError {
			if p.Username == "" || p.Password == "" {
				return nil
			}

			p.user = &models.User{}
			res := db.ResultOrPanic(
				session.Where("username = ?", p.Username).Find(p.user),
			)
			if res.RowsAffected == 0 {
				return api.ErrorWrongCredentials
			}

			decoded := gotils.ResultOrPanic(base64.StdEncoding.DecodeString(p.user.PasswordHash))
			cryptErr := bcrypt.CompareHashAndPassword(decoded, []byte(p.Password))

			if cryptErr != nil {
				return api.ErrorWrongCredentials
			}

			return nil
		}},
	}
}
