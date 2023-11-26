package main

import (
	"db"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/DAtek/gotils"
)

func main() {
	migrator, err := db.NewMigratorFromEnv()
	if err != nil {
		panic(err)
	}

	switch os.Args[1] {
	case "up":
		if err = migrator.Up(); err != nil && !strings.Contains(err.Error(), "no change") {
			panic(err)
		}
	case "down":
		if err = migrator.Down(); err != nil && !strings.Contains(err.Error(), "no change") {
			panic(err)
		}
	case "version":
		version := uint(gotils.ResultOrPanic(strconv.Atoi(os.Args[2])))
		if err = migrator.Migrate(version); err != nil && !strings.Contains(err.Error(), "no change") {
			panic(err)
		}
	case "force":
		version := gotils.ResultOrPanic(strconv.Atoi(os.Args[2]))
		if err = migrator.Force(version); err != nil && !strings.Contains(err.Error(), "no change") {
			panic(err)
		}
	default:
		panic("Choose 'up' | 'down' | 'version' | 'force'")
	}

	fmt.Println("OK")
}
