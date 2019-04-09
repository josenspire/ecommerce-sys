package utils

import "time"

func GenerateNowDateString() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func GenerateAndAddDurationFromNow(duration time.Duration) time.Time {
	return time.Now().Add(duration)
}
