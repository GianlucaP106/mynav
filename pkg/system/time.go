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
	duration := time.Since(t)
	switch {
	case duration < time.Minute:
		seconds := int(duration.Seconds())
		if seconds == 0 {
			return "now"
		}

		if seconds == 1 {
			return "1 second ago"
		}
		return fmt.Sprintf("%d seconds ago", seconds)
	case duration < time.Hour:
		minutes := int(duration.Minutes())
		if minutes == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", minutes)
	case duration < 24*time.Hour:
		hours := int(duration.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	case duration < 30*24*time.Hour:
		days := int(duration.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	default:
		months := int(duration.Hours() / (24 * 30))
		if months == 1 {
			return "1 month ago"
		}
		return fmt.Sprintf("%d months ago", months)
	}
}

// func main() {
// 	// Example timestamp
// 	timestamp := time.Now().Add(-45 * time.Second) // 45 seconds ago
// 	fmt.Println(timeAgo(timestamp))
//
// 	timestamp = time.Now().Add(-3 * time.Hour) // 3 hours ago
// 	fmt.Println(timeAgo(timestamp))
//
// 	timestamp = time.Now().Add(-7 * 24 * time.Hour) // 7 days ago
// 	fmt.Println(timeAgo(timestamp))
// }
