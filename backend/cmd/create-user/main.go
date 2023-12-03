package main

import (
	"api/components/user"
	"db"
	"flag"
	"fmt"

	"github.com/DAtek/gotils"
)

func main() {
	username := flag.String("u", "", "username")
	password := flag.String("p", "", "password")
	flag.Parse()
	session := gotils.ResultOrPanic(db.NewDbFromEnv())
	userStore := user.NewUserStore(session)
	user := gotils.ResultOrPanic(userStore.CreateUser(*username, *password))
	fmt.Printf("New user's ID: %d\n", user.Id)
}
