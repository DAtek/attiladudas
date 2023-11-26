package components

import (
	"gorm.io/gorm"
)

// useful singletons used by tests

var db *gorm.DB = nil

func getDb() *gorm.DB {
	if db == nil {
		db = CreateOrPanic(NewDbFromEnv)
	}

	return db.Begin()
}

var emptyDb *gorm.DB = nil

func getEmptyDb() *gorm.DB {
	if emptyDb == nil {
		emptyDb = CreateOrPanic(newEmptyDb)
	}

	return emptyDb
}

func newEmptyDb() (*gorm.DB, error) {
	return NewDb(
		EnvDbHost.Load(),
		EnvDbUser.Load(),
		EnvDbPassword.Load(),
		EnvDbName.Load(),
		EnvEmptyPostgresPort.Load(),
	)
}

func CreateOrPanic[T interface{}](f func() (T, error)) T {
	result, err := f()
	if err != nil {
		panic(err)
	}
	return result
}
