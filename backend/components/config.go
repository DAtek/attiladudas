package components

import (
	"os"
)

type EnvVariable string

func (v EnvVariable) Load() string {
	return os.Getenv(string(v))
}

const (
	EnvPrivateKey        = EnvVariable("API_PRIVATE_KEY")
	EnvPublicKey         = EnvVariable("API_PUBLIC_KEY")
	EnvDbHost            = EnvVariable("API_DB_HOST")
	EnvDbUser            = EnvVariable("API_DB_USER")
	EnvDbPassword        = EnvVariable("API_DB_PASSWORD")
	EnvDbName            = EnvVariable("API_DB_NAME")
	EnvDbPort            = EnvVariable("API_DB_PORT")
	EnvMediaDir          = EnvVariable("API_MEDIA_DIR")
	EnvMigrationsDir     = EnvVariable("API_MIGRATIONS_DIR")
	EnvSocketPath        = EnvVariable("API_SOCKET_PATH")
	EnvEmptyPostgresPort = EnvVariable("EMPTY_POSTGRES_PORT")
	EnvTestFilesDir      = EnvVariable("TEST_FILES_DIR")
)
