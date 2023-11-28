package token_post

import (
	"api/components/auth"
	"bytes"
	"db"
	"encoding/base64"
	"encoding/json"
	"fibertools"
	"io"
	"net/http"
	"testing"

	"github.com/DAtek/gotils"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestPostToken(t *testing.T) {
	username := "user1"
	password := "password1"
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	hashedPassword := base64.StdEncoding.EncodeToString(passwordHash)

	t.Run("OK with proper credentials", func(t *testing.T) {
		transaction := db.GetTestTransaction()
		defer transaction.Rollback()
		populator := db.NewPopulator(transaction)
		populator.User(map[string]any{"Username": username, "PasswordHash": hashedPassword})

		jwtContext := &auth.MockJwtContext{
			Encode_: func(c *auth.Claims) (string, error) {
				return "LOL", nil
			},
		}

		app := fibertools.NewApp(
			PluginTokenPost(transaction, jwtContext),
		)

		body := gotils.ResultOrPanic(
			json.Marshal(&postTokenBody{Username: username, Password: password}),
		)
		req := gotils.ResultOrPanic(
			http.NewRequest(
				"POST",
				"/api/token/",
				bytes.NewBuffer(body),
			),
		)
		req.Header.Add("Content-Type", "application/json")

		response := gotils.ResultOrPanic(
			app.Test(req),
		)

		respBody := &bytes.Buffer{}
		io.Copy(respBody, response.Body)
		tokenResponse := &tokenResponse{}
		json.Unmarshal(respBody.Bytes(), tokenResponse)

		contentTypeHeader := response.Header.Get("Content-Type")

		assert.Equal(t, http.StatusCreated, response.StatusCode)
		assert.Equal(t, "application/json", contentTypeHeader)
		assert.NotEqual(t, "", tokenResponse.Token)
	})

	badRequestScenarios := []struct {
		name     string
		username string
		password string
	}{
		{"Wrong username1", "", password},
		{"Wrong username2", "asd", password},
		{"Wrong password1", username, ""},
		{"Wrong password2", username, "asd"},
	}
	for _, scenario := range badRequestScenarios {
		t.Run(scenario.name, func(t *testing.T) {
			transaction := db.GetTestTransaction()
			defer transaction.Rollback()
			populator := db.NewPopulator(transaction)
			populator.User(map[string]any{"Username": username, "PasswordHash": hashedPassword})
			app := fibertools.NewApp(
				PluginTokenPost(transaction, &auth.MockJwtContext{}),
			)
			body, _ := json.Marshal(postTokenBody{Username: scenario.username, Password: scenario.password})

			req := gotils.ResultOrPanic(
				http.NewRequest(
					"POST",
					"/api/token/",
					bytes.NewBuffer(body),
				),
			)
			req.Header.Add("Content-Type", "application/json")

			resp := gotils.ResultOrPanic(app.Test(req))

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	}

	t.Run("Bad request if body isn't json", func(t *testing.T) {
		transaction := db.GetTestTransaction()
		defer transaction.Rollback()

		app := fibertools.NewApp(
			PluginTokenPost(transaction, &auth.MockJwtContext{}),
		)

		body := gotils.ResultOrPanic(
			json.Marshal(&postTokenBody{Username: username, Password: password}),
		)
		req := gotils.ResultOrPanic(
			http.NewRequest(
				"POST",
				"/api/token/",
				bytes.NewBuffer(body),
			),
		)
		req.Header.Add("Content-Type", "application/json")

		response := gotils.ResultOrPanic(
			app.Test(req),
		)

		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	})
}
