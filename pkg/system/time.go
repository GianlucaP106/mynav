package system

import (
	"fmt"
	"strconv"
	"time"
)

func TimeFormat() string {
	return "02 Jan 15:04"
}

func UnixTime(t string) time.Time {
	activityInt, err := strconv.Atoi(t)
	if err != nil {
	}

	time := time.Unix(int64(activityInt), 0)
	return time
}

func TimeAgo(t time.Time) string {
	timeStr := TimeDeltaStr(time.Since(t))
	if timeStr == "0 seconds" {
		return "now"
	}
	return timeStr + " ago"
}

func TimeDeltaStr(duration time.Duration) string {
	switch {
	case duration < time.Minute:
		seconds := int(duration.Seconds())
		if seconds == 0 {
			return "0 seconds"
		}

		if seconds == 1 {
			return "1 second"
		}
		return fmt.Sprintf("%d seconds", seconds)
	case duration < time.Hour:
		minutes := int(duration.Minutes())
		if minutes == 1 {
			return "1 minute"
		}
		return fmt.Sprintf("%d minutes", minutes)
	case duration < 24*time.Hour:
		hours := int(duration.Hours())
		if hours == 1 {
			return "1 hour"
		}
		return fmt.Sprintf("%d hours", hours)
	case duration < 30*24*time.Hour:
		days := int(duration.Hours() / 24)
		if days == 1 {
			return "1 day"
		}
		return fmt.Sprintf("%d days", days)
	default:
		months := int(duration.Hours() / (24 * 30))
		if months == 1 {
			return "1 month"
		}
		return fmt.Sprintf("%d months", months)
	}
}
