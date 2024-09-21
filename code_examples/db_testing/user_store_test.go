package db_test

import (
	"db_test/db"
	"testing"

	"github.com/DAtek/gotils"
	"github.com/stretchr/testify/assert"

	"gorm.io/gorm"
)

func TestCreateUser(t *testing.T) {
	t.Run(
		"Creates user",
		db.TestWithTransaction(func(t *testing.T, tx *gorm.DB) {
			// given
			name := "Asimov"

			// when
			user := gotils.ResultOrPanic(CreateUser(tx, name))

			// then
			assert.Equal(t, name, user.Name)
			savedUser := db.User{}
			tx.Find(&savedUser).Where("id = ?", user.Id)
			assert.Equal(t, name, savedUser.Name)
		}))

	t.Run(
		"Returns error if user with the same name already exists",
		db.TestWithTransaction(func(t *testing.T, tx *gorm.DB) {
			// given
			name := "Asimov"
			CreateUser(tx, name)

			// when
			_, err := CreateUser(tx, name)

			// then
			assert.Error(t, err)
		}))
}

func TestDeleteUser(t *testing.T) {
	t.Run(
		"Deletes user",
		db.TestWithTransaction(func(t *testing.T, tx *gorm.DB) {
			// given
			user := db.User{Name: "Isaac"}
			res := tx.Create(&user)
			if res.Error != nil {
				panic(res.Error)
			}

			// when
			err := DeleteUser(tx, user.Id)

			// then
			assert.Nil(t, err)
			count := int64(0)
			tx.Find(&db.User{}).Where("id = ?", user.Id).Count(&count)
			assert.Equal(t, int64(0), count)
		}))

	t.Run(
		"Returns error if user not exists",
		db.TestWithTransaction(func(t *testing.T, tx *gorm.DB) {
			// when
			err := DeleteUser(tx, 1)

			// then
			assert.Error(t, err)
		}))

	t.Run(
		"Returns error if unexpected event happens",
		db.TestWithTransaction(func(t *testing.T, tx *gorm.DB) {
			// given
			user := db.User{Name: "Isaac"}
			res := tx.Create(&user)
			if res.Error != nil {
				panic(res.Error)
			}

			tx.Rollback()

			// when
			err := DeleteUser(tx, user.Id)

			// then
			assert.Error(t, err)
		}))
}
