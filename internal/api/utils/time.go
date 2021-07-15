package utils

import "time"

const DefaultLayoutDB = "2006-01-02 15:04:05"

type GetTimeNowFn func() time.Time

func GetTimeNow() time.Time {
	return time.Now()
}
