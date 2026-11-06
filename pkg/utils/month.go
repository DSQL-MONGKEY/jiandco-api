package utils

import (
	"fmt"
	"time"
)

const MonthLayout = "2006-01"

func ParseMonth(value string, location *time.Location) (start time.Time, end time.Time, err error) {
	now := time.Now().In(location)
	if value == "" {
		start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, location)
	} else {
		parsed, parseErr := time.ParseInLocation(MonthLayout, value, location)
		if parseErr != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("month must use YYYY-MM format")
		}
		start = time.Date(parsed.Year(), parsed.Month(), 1, 0, 0, 0, 0, location)
	}

	end = start.AddDate(0, 1, 0)
	return start, end, nil
}

func MonthString(t time.Time, location *time.Location) string {
	return t.In(location).Format(MonthLayout)
}
