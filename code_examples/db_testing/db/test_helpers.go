package db

import (
	"sync"
	"testing"

	"github.com/DAtek/gotils"
	"gorm.io/gorm"
)

var GetTestConn = sync.OnceValue(func() *gorm.DB {
	return gotils.ResultOrPanic(NewConnFromEnv())
})

func GetTestTransaction() *gorm.DB {
	return GetTestConn().Begin()
}

func TestWithTransaction(
	transactionalTest func(t *testing.T, tx *gorm.DB),
) func(t *testing.T) {
	return func(t *testing.T) {
		tx := GetTestTransaction()
		t.Cleanup(func() { tx.Rollback() })
		transactionalTest(t, tx)
	}
}
