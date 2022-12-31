package component_mocks

import "attiladudas/backend/models"

type MockUserStore struct {
	GetUser_    func(string) (*models.User, error)
	CreateUser_ func(string, string) (*models.User, error)
}

func (db *MockUserStore) GetUser(username string) (*models.User, error) {
	return db.GetUser_(username)
}

func (db *MockUserStore) CreateUser(username, passwordHash string) (*models.User, error) {
	return db.CreateUser_(username, passwordHash)
}
