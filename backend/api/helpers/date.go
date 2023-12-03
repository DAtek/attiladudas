package helpers

import (
	"time"

	"gorm.io/datatypes"
)

func DateToISO8601(date *datatypes.Date) string {
	return time.Time(*date).Format(time.DateOnly)
}

func DateFromISO8601(date string) (*datatypes.Date, error) {
	t, err := time.Parse(time.DateOnly, date)
	if err != nil {
		return nil, err
	}
	d := datatypes.Date(t)
	return &d, nil
}
