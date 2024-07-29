package utils

import "time"

func IntToDuration(s int) time.Duration {
	return time.Duration(s) * time.Second
}
