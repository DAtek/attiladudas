package db

import (
	"sync"
	"time"

	"gorm.io/gorm"
)

var _db *gorm.DB

// Use only in tests!
func GetTestTransaction() *gorm.DB {
	if _db == nil {
		_db = createOrPanic(NewDbFromEnv)
	}

	return _db.Begin()
}

func ResultOrPanic(res *gorm.DB) *gorm.DB {
	if res.Error != nil {
		panic(res.Error)
	}

	return res
}

func createOrPanic[T any](f func() (T, error)) T {
	result, err := f()
	if err != nil {
		panic(err)
	}
	return result
}

type txManager struct {
	transaction  *gorm.DB
	mutex        *sync.Mutex
	lockDuration time.Duration
}

// Use only in tests!
func NewTxManager(tx *gorm.DB, lockDuration time.Duration) *txManager {
	return &txManager{
		transaction:  tx,
		mutex:        &sync.Mutex{},
		lockDuration: lockDuration,
	}
}

func (m *txManager) GetTransaction() *gorm.DB {
	m.mutex.Lock()
	go func() {
		// Giving time for using the transaction
		time.Sleep(m.lockDuration)
		defer m.mutex.Unlock()
	}()
	return m.transaction
}
