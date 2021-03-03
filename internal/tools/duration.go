package tools

import (
	"regexp"
	"strconv"
	"time"
)

var durationRegex = regexp.MustCompile(`(?:(\d+)d)?(?:(\d+)h)?(?:(\d+)m)?`)

func ParseDuration(str string) time.Duration {
	matches := durationRegex.FindStringSubmatch(str)

	days := parseInt64(matches[1])
	hours := parseInt64(matches[2])
	minutes := parseInt64(matches[3])

	hour := int64(time.Hour)
	return time.Duration(days*24*hour + hours*hour + minutes*int64(time.Minute))
}

func parseInt64(value string) int64 {
	if len(value) == 0 {
		return 0
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return int64(parsed)
}
