package components

import (
	"attiladudas/backend/models"
	"encoding/base64"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

/*
 * A running postgres instance is required to run these tests
 */
func TestUserStore(t *testing.T) {
	t.Run("CreateUser saves user to the db with proper data", func(t *testing.T) {
		db := getDb()
		defer db.Rollback()
		store := NewUserStore(db)

		username := "Maria"
		password := "Password1"
		returnedUser, err := store.CreateUser(username, password)
		user := &models.User{}
		db.Where("username = ?", username).Find(user)
		decoded, _ := base64.StdEncoding.DecodeString(user.PasswordHash)
		passwordErr := bcrypt.CompareHashAndPassword(decoded, []byte(password))

		assert.Nil(t, err)
		assert.Nil(t, passwordErr)
		assert.Equal(t, returnedUser.Id, user.Id)
		assert.Equal(t, user, returnedUser)

	})

	t.Run("GetUser throws not found error", func(t *testing.T) {
		db := getDb()
		defer db.Rollback()
		store := NewUserStore(db)
		_, err := store.GetUser("Max")

		assert.EqualError(t, NotFoundError, err.Error())
	})

	t.Run("Returns error if something unexpected happens", func(t *testing.T) {
		db := getEmptyDb()
		store := NewUserStore(db)
		_, err := store.GetUser("")

		assert.Error(t, err)
		assert.False(t, strings.Contains(err.Error(), NotFoundError.Error()))
	})

	t.Run("GetUser returns proper user", func(t *testing.T) {
		db := getDb()
		defer db.Rollback()
		store := NewUserStore(db)
		users := []models.User{
			{Username: "John", PasswordHash: "asd"},
			{Username: "Martha", PasswordHash: "dfg"},
		}

		db.Create(&users)

		user, err := store.GetUser(users[1].Username)

		assert.Nil(t, err)
		assert.Equal(t, users[1].Username, user.Username)
		assert.Equal(t, users[1].PasswordHash, user.PasswordHash)
	})
}
