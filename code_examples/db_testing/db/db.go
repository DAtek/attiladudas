package db

import (
	"fmt"

	"github.com/DAtek/env"
	"github.com/DAtek/gotils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func NewConnFromEnv() (*gorm.DB, error) {
	loadConfig := env.NewLoader[EnvConfig]()

	config := gotils.ResultOrPanic(loadConfig())

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=UTC",
		config.AppDbHost,
		config.AppDbUser,
		config.AppDbPassword,
		config.AppDbName,
		config.AppDbPort,
	)

	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		SkipDefaultTransaction: true,
	})
}

func CommitTxFunc(conn *gorm.DB, f func(tx *gorm.DB) error) error {
	return conn.Transaction(f)
}

func MigrateUp(conn *gorm.DB) {
	panicOnError(conn.Exec(`
CREATE TABLE "user" (
	id SERIAL PRIMARY KEY NOT NULL,
	name VARCHAR NOT NULL UNIQUE
);`))
}

func MigrateDown(conn *gorm.DB) {
	panicOnError(conn.Exec(`DROP TABLE "user";`))
}

func panicOnError(result *gorm.DB) {
	if result.Error != nil {
		panic(result.Error)
	}
}
