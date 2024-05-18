package utils

import (
	"runtime"
	"time"
)

type OS = uint

const (
	Darwin OS = iota
	Linux
	Unsuported
)

func DetectOS() OS {
	switch runtime.GOOS {
	case "darwin":
		return Darwin
	case "linux":
		return Linux
	default:
		return Unsuported
	}
}

func IsBeforeOneHourAgo(timestamp time.Time) bool {
	now := time.Now()
	oneHourAgo := now.Add(-1 * time.Hour)
	return timestamp.Before(oneHourAgo)
}
