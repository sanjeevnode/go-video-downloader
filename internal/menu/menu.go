package menu

import (
	"bufio"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/sanjeevnode/go-video-downloader/internal/downloader"
	"github.com/sanjeevnode/go-video-downloader/internal/search"
	"github.com/sanjeevnode/go-video-downloader/internal/utils"
)

func ShowMainMenu() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\nSelect an option:")
		fmt.Println("1. Search YouTube")
		fmt.Println("2. Download from URL")
		fmt.Println("3. Exit")

		fmt.Print("Enter choice: ")
		choiceStr, _ := reader.ReadString('\n')
		choiceStr = strings.TrimSpace(choiceStr)
		choice, err := strconv.Atoi(choiceStr)
		if err != nil {
			fmt.Println("Invalid input, try again.")
			continue
		}

		switch choice {
		case 1:
			handleSearch(reader)
		case 2:
			handleDownloadFromURL(reader)
		case 3:
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Please enter 1, 2, or 3.")
		}
	}
}

func handleSearch(reader *bufio.Reader) {
	fmt.Print("Enter search keyword: ")
	keyword, _ := reader.ReadString('\n')
	keyword = strings.TrimSpace(keyword)

	videos, err := search.SearchYouTube(keyword, 10)
	if err != nil {
		fmt.Println("Error searching YouTube:", err)
		return
	}

	// for i, v := range videos {
	// 	fmt.Printf("%d. Title: %s ChannelName: %s Views: %s Duration: %s (Published: %s)\n", i+1, v.Title, v.ChannelName, v.ViewCount, v.Duration, v.PublishedAt)
	// }
	utils.PrintVideosTable(videos)

	fmt.Println("Enter video number to download or 0 to return to main menu:")
	choiceStr, _ := reader.ReadString('\n')
	choiceStr = strings.TrimSpace(choiceStr)
	choice, err := strconv.Atoi(choiceStr)
	if err != nil || choice < 0 || choice > len(videos) {
		fmt.Println("Invalid choice.")
		return
	}

	if choice == 0 {
		return
	}

	selectedVideo := videos[choice-1]
	fmt.Println("You selected:", selectedVideo.Title)

	url := "https://www.youtube.com/watch?v=" + selectedVideo.VideoID
	fmt.Println("Video URL:", url)

	format, err1 := handleFormatSelection(reader)
	if err1 != nil {
		fmt.Println("Error selecting format:", err1)
		return
	}

	downError := downloader.Download(url, format)
	if downError != nil {
		fmt.Println("Error downloading:", downError)
		return
	}

}

func handleDownloadFromURL(reader *bufio.Reader) {
	fmt.Print("Paste YouTube URL: ")
	url, _ := reader.ReadString('\n')
	url = strings.TrimSpace(url)
	errUrl := ValidateYouTubeURL(url)
	if errUrl != nil {
		fmt.Println("Invalid YouTube URL:", errUrl)
		return
	}

	format, err1 := handleFormatSelection(reader)
	if err1 != nil {
		fmt.Println("Error selecting format:", err1)
		return
	}

	err := downloader.Download(url, format)
	if err != nil {
		fmt.Println("Error downloading:", err)
		return
	}
}

func handleFormatSelection(reader *bufio.Reader) (int, error) {
	fmt.Println("Choose format to download:")
	fmt.Println("1. MP3 (audio only)")
	fmt.Println("2. 480p")
	fmt.Println("3. 720p")
	fmt.Println("4. 1080p")

	formatStr, _ := reader.ReadString('\n')
	formatStr = strings.TrimSpace(formatStr)
	format, err := strconv.Atoi(formatStr)
	if err != nil || format < 1 || format > 4 {
		fmt.Println("Invalid format choice.")
		return -1, err
	}
	return format, nil
}

func ValidateYouTubeURL(videoURL string) error {
	parsed, err := url.Parse(videoURL)
	if err != nil {
		return errors.New("invalid URL format")
	}

	host := parsed.Hostname()
	if !(strings.Contains(host, "youtube.com") || strings.Contains(host, "youtu.be")) {
		return errors.New("URL is not a valid YouTube link")
	}

	// For youtube.com URLs, check if there's a "v" query parameter
	if strings.Contains(host, "youtube.com") {
		vals := parsed.Query()
		if vals.Get("v") == "" {
			return errors.New("missing video ID in URL query parameter")
		}
	}

	// For youtu.be URLs, path should not be empty (it contains video ID)
	if strings.Contains(host, "youtu.be") {
		if len(strings.Trim(parsed.Path, "/")) == 0 {
			return errors.New("missing video ID in URL path")
		}
	}

	return nil
}
