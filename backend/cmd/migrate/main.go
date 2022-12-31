package main

import (
	"attiladudas/backend/components"
	"fmt"
	"os"
	"strings"
)

func main() {
	migrator, err := components.NewMigratorFromEnv()
	if err != nil {
		panic(err)
	}

	switch os.Args[1] {
	case "up":
		if err = migrator.Up(); err != nil {
			if !strings.Contains(err.Error(), "no change") {
				panic(err)
			}
		}
	case "down":
		if err = migrator.Down(); err != nil {
			if !strings.Contains(err.Error(), "no change") {
				panic(err)
			}
		}
	default:
		panic("Choose 'up' or 'down'")
	}

	fmt.Println("Migrations ran successfully.")
}
