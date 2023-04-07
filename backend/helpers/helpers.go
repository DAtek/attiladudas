package helpers

import (
	"time"

	"gorm.io/datatypes"
)

func DateFromISO8601(isoDateStr string) (*datatypes.Date, error) {
	time, err := time.Parse("2006-01-02", isoDateStr)
	if err != nil {
		return nil, err
	}
	date := datatypes.Date(time)
	return &date, nil
}

func DateToISO8601(date datatypes.Date) string {
	rfcTime := time.Time(date).Format(time.RFC3339)
	return rfcTime[0:10]
}

func DateFromISO8601Panic(isoDateStr string) *datatypes.Date {
	date, err := DateFromISO8601(isoDateStr)
	if err != nil {
		panic(err)
	}
	return date
}
