package components

import (
	"attiladudas/backend/models"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type IUserStore interface {
	GetUser(username string) (*models.User, error)
	CreateUser(username, password string) (*models.User, error)
}

type userStore struct {
	db *gorm.DB
}

func NewUserStore(db *gorm.DB) IUserStore {
	return &userStore{db}
}

func (store *userStore) GetUser(username string) (*models.User, error) {
	user := &models.User{}
	result := store.db.Where("username = ?", username).Find(user)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, NotFoundError
	}

	return user, nil
}

func (store *userStore) CreateUser(username, password string) (*models.User, error) {
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	hashedPassword := base64.StdEncoding.EncodeToString(passwordHash)
	user := &models.User{Username: username, PasswordHash: hashedPassword}
	result := store.db.Create(&user)
	return user, result.Error
}
