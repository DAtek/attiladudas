package db_test

import (
	"db_test/db"
	"errors"

	"gorm.io/gorm"
)

func CreateUser(conn *gorm.DB, name string) (db.User, error) {
	user := db.User{Name: name}
	res := conn.Create(&user)
	return user, res.Error
}

func DeleteUser(conn *gorm.DB, id int) error {
	res := conn.Where("id = ?", id).Delete(&db.User{})

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return errors.New("USER_NOT_EXISTS")
	}

	return nil
}
