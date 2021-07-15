package utils

import "time"

const DefaultLayoutDB = "2006-01-02 15:04:05"

func GetTimeNow() time.Time {
	return time.Now()
}
