package main

import (
	"db_test"
	"db_test/db"
	"os"
	"strconv"

	"github.com/DAtek/gotils"
	"gorm.io/gorm"
)

func main() {
	conn := gotils.ResultOrPanic(db.NewConnFromEnv())
	id := gotils.ResultOrPanic(strconv.Atoi(os.Args[1]))

	gotils.NilOrPanic(
		db.CommitTxFunc(
			conn,
			func(tx *gorm.DB) error {
				return db_test.DeleteUser(tx, id)
			},
		),
	)
}
