package main

import (
	"db_test"
	"db_test/db"
	"fmt"
	"os"

	"github.com/DAtek/gotils"
	"gorm.io/gorm"
)

func main() {
	conn := gotils.ResultOrPanic(db.NewConnFromEnv())
	name := os.Args[1]

	gotils.NilOrPanic(
		db.CommitTxFunc(
			conn,
			func(tx *gorm.DB) error {
				user, err := db_test.CreateUser(tx, name)
				fmt.Printf("%d\n", user.Id)
				return err
			},
		),
	)
}
