package db

import "github.com/DAtek/gotils"

const (
	EnvDbUser        = gotils.EnvConfig("POSTGRES_USER")
	EnvDbPassword    = gotils.EnvConfig("POSTGRES_PASSWORD")
	EnvDbHost        = gotils.EnvConfig("POSTGRES_HOST")
	EnvDbPort        = gotils.EnvConfig("POSTGRES_PORT")
	EnvDbName        = gotils.EnvConfig("POSTGRES_DB")
	EnvMigrationsDir = gotils.EnvConfig("MIGRATIONS_DIR")
)
