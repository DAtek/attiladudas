package helpers

import (
	"time"

	"gorm.io/datatypes"
)

func DateToISO8601(date *datatypes.Date) string {
	return time.Time(*date).Format(time.DateOnly)
}
