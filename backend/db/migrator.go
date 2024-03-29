package db

import (
	"embed"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

//go:embed migrations/*
var migrationsDir embed.FS

func NewMigratorFromEnv() (*migrate.Migrate, error) {
	return migrate.New(
		"file://"+EnvMigrationsDir.Load(),
		fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			EnvDbUser.Load(),
			EnvDbPassword.Load(),
			EnvDbHost.Load(),
			EnvDbPort.Load(),
			EnvDbName.Load(),
		),
	)
}
