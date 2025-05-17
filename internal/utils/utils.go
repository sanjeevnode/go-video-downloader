package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/sanjeevnode/go-video-downloader/internal/video"
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

func wrapText(s string, width int) []string {
	var lines []string
	for len(s) > width {
		// Find last space within width to break nicely
		splitAt := strings.LastIndex(s[:width], " ")
		if splitAt == -1 {
			splitAt = width // no space, just hard break
		}
		lines = append(lines, s[:splitAt])
		s = s[splitAt:]
		s = strings.TrimLeft(s, " ")
	}
	if len(s) > 0 {
		lines = append(lines, s)
	}
	return lines
}

func PrintVideosTable(videos []video.Video) {
	// Define column widths
	const (
		idxWidth       = 3
		titleWidth     = 50
		channelWidth   = 20
		viewsWidth     = 10
		durationWidth  = 8
		publishedWidth = 12
	)

	// Print header
	fmt.Printf("%-*s %-*s %-*s %-*s %-*s %-*s\n",
		idxWidth, "No.",
		titleWidth, "Title",
		channelWidth, "Channel",
		viewsWidth, "Views",
		durationWidth, "Duration",
		publishedWidth, "Published")

	// Print separator
	fmt.Printf("%s\n", strings.Repeat("-", idxWidth+titleWidth+channelWidth+viewsWidth+durationWidth+publishedWidth+6))

	for i, v := range videos {
		// Wrap title and channel name
		titleLines := wrapText(v.Title, titleWidth)
		channelLines := wrapText(v.ChannelName, channelWidth)

		maxLines := len(titleLines)
		if len(channelLines) > maxLines {
			maxLines = len(channelLines)
		}

		for line := 0; line < maxLines; line++ {
			idxStr := ""
			viewsStr := ""
			durationStr := ""
			publishedStr := ""

			if line == 0 {
				idxStr = fmt.Sprintf("%d.", i+1)
				viewsStr = v.ViewCount
				durationStr = v.Duration
				publishedStr = v.PublishedAt
			}

			t := ""
			if line < len(titleLines) {
				t = titleLines[line]
			}

			c := ""
			if line < len(channelLines) {
				c = channelLines[line]
			}

			fmt.Printf("%-*s %-*s %-*s %-*s %-*s %-*s\n",
				idxWidth, idxStr,
				titleWidth, t,
				channelWidth, c,
				viewsWidth, viewsStr,
				durationWidth, durationStr,
				publishedWidth, publishedStr)
		}
	}
}
