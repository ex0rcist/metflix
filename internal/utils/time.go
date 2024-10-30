package utils

import "time"

// Convert int (seconds) to time.Duration
func IntToDuration(s int) time.Duration {
	return time.Duration(s) * time.Second
}
