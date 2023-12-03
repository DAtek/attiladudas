package user

import (
	"db"
	"db/models"
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestUserStore(t *testing.T) {
	t.Run("CreateUser saves user to the db with proper data", func(t *testing.T) {
		transaction := db.GetTestTransaction()
		defer transaction.Rollback()
		store := NewUserStore(transaction)

		username := "Maria"
		password := "Password1"
		returnedUser, err := store.CreateUser(username, password)
		user := &models.User{}
		transaction.Where("username = ?", username).Find(user)
		decoded, _ := base64.StdEncoding.DecodeString(user.PasswordHash)
		passwordErr := bcrypt.CompareHashAndPassword(decoded, []byte(password))

		assert.Nil(t, err)
		assert.Nil(t, passwordErr)
		assert.Equal(t, returnedUser.Id, user.Id)
		assert.Equal(t, user, returnedUser)

	})
}
