package api

import (
	"github.com/DAtek/gotils"
)

const (
	EnvPrivateKey    = gotils.EnvConfig("API_PRIVATE_KEY")
	EnvPublicKey     = gotils.EnvConfig("API_PUBLIC_KEY")
	EnvMediaDir      = gotils.EnvConfig("API_MEDIA_DIR")
	EnvMigrationsDir = gotils.EnvConfig("API_MIGRATIONS_DIR")
	EnvSocketPath    = gotils.EnvConfig("API_SOCKET_PATH")
)
