package search

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/sanjeevnode/go-video-downloader/internal/config"
	"github.com/sanjeevnode/go-video-downloader/internal/utils"
	"github.com/sanjeevnode/go-video-downloader/internal/video"
)

func SearchYouTube(query string, maxResults int) ([]video.Video, error) {
	apiKey := config.GetAPIKey()
	searchURL := "https://www.googleapis.com/youtube/v3/search"

	params := url.Values{}
	params.Set("part", "snippet")
	params.Set("q", query)
	params.Set("maxResults", fmt.Sprintf("%d", maxResults))
	params.Set("type", "video")
	params.Set("key", apiKey)

	fullURL := searchURL + "?" + params.Encode()
	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// fmt.Println("Response Body (Search):", string(bodyBytes))

	var searchData struct {
		Items []struct {
			ID struct {
				VideoID string `json:"videoId"`
			} `json:"id"`
			Snippet struct {
				Title       string `json:"title"`
				PublishedAt string `json:"publishedAt"`
			} `json:"snippet"`
		} `json:"items"`
	}

	if err := json.Unmarshal(bodyBytes, &searchData); err != nil {
		return nil, err
	}

	// Collect all video IDs for the next API call
	var videoIDs []string
	for _, item := range searchData.Items {
		videoIDs = append(videoIDs, item.ID.VideoID)
	}

	if len(videoIDs) == 0 {
		return nil, nil
	}

	// Call videos endpoint to get extra info
	videosURL := "https://www.googleapis.com/youtube/v3/videos"
	videoParams := url.Values{}
	videoParams.Set("part", "snippet,contentDetails,statistics")
	videoParams.Set("id", strings.Join(videoIDs, ","))
	videoParams.Set("key", apiKey)

	fullVideosURL := videosURL + "?" + videoParams.Encode()
	resp2, err := http.Get(fullVideosURL)
	if err != nil {
		return nil, err
	}
	defer resp2.Body.Close()

	bodyBytes2, err := io.ReadAll(resp2.Body)
	if err != nil {
		return nil, err
	}
	// fmt.Println("Response Body (Videos):", string(bodyBytes2))

	var videosData struct {
		Items []struct {
			ID      string `json:"id"`
			Snippet struct {
				ChannelTitle string `json:"channelTitle"`
			} `json:"snippet"`
			ContentDetails struct {
				Duration string `json:"duration"` // ISO 8601 duration format, e.g. PT15M33S
			} `json:"contentDetails"`
			Statistics struct {
				ViewCount string `json:"viewCount"`
			} `json:"statistics"`
		} `json:"items"`
	}

	if err := json.Unmarshal(bodyBytes2, &videosData); err != nil {
		return nil, err
	}

	// Map videoID to video details for quick lookup
	videoDetailsMap := make(map[string]struct {
		Duration    string
		ChannelName string
		ViewCount   string
	})
	for _, item := range videosData.Items {
		videoDetailsMap[item.ID] = struct {
			Duration    string
			ChannelName string
			ViewCount   string
		}{
			Duration:    item.ContentDetails.Duration,
			ChannelName: item.Snippet.ChannelTitle,
			ViewCount:   item.Statistics.ViewCount,
		}
	}

	var results []video.Video
	for _, item := range searchData.Items {
		details := videoDetailsMap[item.ID.VideoID]
		results = append(results, video.Video{
			Title:       item.Snippet.Title,
			VideoID:     item.ID.VideoID,
			PublishedAt: item.Snippet.PublishedAt,
			Duration:    utils.ParseYouTubeDuration(details.Duration),
			ChannelName: details.ChannelName,
			ViewCount:   details.ViewCount,
		})
	}

	return results, nil
}
