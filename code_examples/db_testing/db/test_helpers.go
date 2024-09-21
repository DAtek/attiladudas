package db

import (
	"testing"

	"github.com/DAtek/gotils"
	"gorm.io/gorm"
)

var testConn *gorm.DB

func GetTestConn() *gorm.DB {
	if testConn == nil {
		testConn = gotils.ResultOrPanic(NewConnFromEnv())
	}

	return testConn
}

func GetTestTransaction() *gorm.DB {
	return GetTestConn().Begin()
}

func TestWithTransaction(
	transactionalTest func(t *testing.T, tx *gorm.DB),
) func(t *testing.T) {
	tx := GetTestTransaction()

	return func(t *testing.T) {
		defer tx.Rollback()
		transactionalTest(t, tx)
	}
}
