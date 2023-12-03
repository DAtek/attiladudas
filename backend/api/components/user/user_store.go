package user

import (
	"db/models"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type IUserStore interface {
	CreateUser(username, password string) (*models.User, error)
}

type userStore struct {
	db *gorm.DB
}

func NewUserStore(db *gorm.DB) IUserStore {
	return &userStore{db}
}

func (store *userStore) CreateUser(username, password string) (*models.User, error) {
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	hashedPassword := base64.StdEncoding.EncodeToString(passwordHash)
	user := &models.User{Username: username, PasswordHash: hashedPassword}
	result := store.db.Create(&user)
	return user, result.Error
}
