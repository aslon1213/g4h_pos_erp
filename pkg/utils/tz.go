package utils

import "time"

func GetTimeZone() *time.Location {
	loc, err := time.LoadLocation("Asia/Tashkent")
	if err != nil {
		return time.UTC // default to UTC if error
	}
	return loc
}
