package main

import (
	"attiladudas/backend/components"
	"fmt"
	"os"
)

func main() {
	username := os.Args[1]
	password := os.Args[2]
	db, err := components.NewDbFromEnv()
	userStore := components.NewUserStore(db)

	if err != nil {
		panic(err)
	}

	user, createErr := userStore.CreateUser(username, password)

	if err != createErr {
		panic(err)
	}

	fmt.Printf("New user's ID: %d\n", user.Id)
}
