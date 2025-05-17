package utils

import (
	"fmt"
	"regexp"
	"strconv"
)

func ParseYouTubeDuration(isoDuration string) string {
	// Example isoDuration: "PT1H2M10S", "PT15M33S", "PT45S", "PT2H"
	re := regexp.MustCompile(`PT(?:(\d+)H)?(?:(\d+)M)?(?:(\d+)S)?`)
	matches := re.FindStringSubmatch(isoDuration)

	if len(matches) == 0 {
		return "00:00:00"
	}

	hours := 0
	minutes := 0
	seconds := 0

	if matches[1] != "" {
		hours, _ = strconv.Atoi(matches[1])
	}
	if matches[2] != "" {
		minutes, _ = strconv.Atoi(matches[2])
	}
	if matches[3] != "" {
		seconds, _ = strconv.Atoi(matches[3])
	}

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}
