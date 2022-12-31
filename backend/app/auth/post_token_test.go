package app_auth

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"attiladudas/backend/components"
	"attiladudas/backend/models"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestPostToken(t *testing.T) {
	username := "user1"
	password := "password1"
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	hashedPassword := base64.StdEncoding.EncodeToString(passwordHash)

	t.Run("OK with proper credentials", func(t *testing.T) {
		app := newMockApp()
		body, _ := json.Marshal(&postTokenBody{Username: username, Password: password})
		app.userStore.GetUser_ = func(s string) (*models.User, error) {
			return &models.User{
					Id:           1,
					Username:     username,
					PasswordHash: hashedPassword,
				},
				nil
		}
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(
			"POST",
			"/api/token/",
			bytes.NewBuffer(body),
		)

		app.engine.ServeHTTP(w, req)
		response := tokenResponse{}
		decoder := json.NewDecoder(w.Body)
		err := decoder.Decode(&response)
		contentTypeHeader := w.Header().Get("content-type")

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Equal(t, "application/json; charset=utf-8", contentTypeHeader)
		assert.Nil(t, err)
	})

	badRequestScenarios := []struct {
		name     string
		username string
		password string
	}{
		{"Wrong password", username, ""},
		{"Wrong username", "", password},
	}
	for _, scenario := range badRequestScenarios {
		t.Run(scenario.name, func(t *testing.T) {
			app := newMockApp()
			app.userStore.GetUser_ = func(s string) (*models.User, error) {
				if s == username {
					return &models.User{
						Username:     username,
						PasswordHash: hashedPassword,
						Id:           1,
					}, nil
				}

				return nil, components.NotFoundError
			}
			body, _ := json.Marshal(postTokenBody{Username: scenario.username, Password: scenario.password})
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(
				"POST",
				"/api/token/",
				bytes.NewBuffer(body),
			)

			app.engine.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)
		})
	}

	t.Run("Internal server error if can't base64decode the password", func(t *testing.T) {
		app := newMockApp()
		app.userStore.GetUser_ = func(s string) (*models.User, error) {
			return &models.User{
				Username:     username,
				PasswordHash: "#&@#đ[/=%",
			}, nil
		}

		w := httptest.NewRecorder()
		body, _ := json.Marshal(postTokenBody{Username: username, Password: password})
		req, _ := http.NewRequest(
			"POST",
			"/api/token/",
			bytes.NewBuffer(body),
		)
		app.engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Internal server error if unexpected happens", func(t *testing.T) {
		app := newMockApp()
		app.userStore.GetUser_ = func(s string) (*models.User, error) {
			return nil, errors.New("BLACK DEATH")
		}

		w := httptest.NewRecorder()
		body, _ := json.Marshal(postTokenBody{Username: username, Password: password})
		req, _ := http.NewRequest(
			"POST",
			"/api/token/",
			bytes.NewBuffer(body),
		)
		app.engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Bad request if body isn't json", func(t *testing.T) {
		app := newMockApp()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(
			"POST",
			"/api/token/",
			bytes.NewBuffer([]byte("username=Mozart;password=asd")),
		)
		app.engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
