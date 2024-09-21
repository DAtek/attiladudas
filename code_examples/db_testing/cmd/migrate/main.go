package main

import (
	"db_test/db"
	"os"

	"github.com/DAtek/gotils"
	"gorm.io/gorm"
)

func main() {
	conn := gotils.ResultOrPanic(db.NewConnFromEnv())
	command, found := commands[os.Args[1]]

	if !found {
		panic("Invalid arg. Allowed: up | down")
	}

	command(conn)
}

var commands = map[string]func(conn *gorm.DB){
	"up":   db.MigrateUp,
	"down": db.MigrateDown,
}
