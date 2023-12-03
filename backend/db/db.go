package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func NewDbFromEnv() (*gorm.DB, error) {
	return NewDb(
		EnvDbHost.Load(),
		EnvDbUser.Load(),
		EnvDbPassword.Load(),
		EnvDbName.Load(),
		EnvDbPort.Load(),
	)
}

func NewDb(host, user, password, dbName, port string) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		host,
		user,
		password,
		dbName,
		port,
	)
	db, dbErr := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		SkipDefaultTransaction: true,
	})

	if dbErr != nil {
		return nil, dbErr
	}

	return db, nil
}
