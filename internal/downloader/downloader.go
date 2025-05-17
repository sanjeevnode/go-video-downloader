package downloader

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

//go:embed yt-dlp.exe ffmpeg.exe ffprobe.exe
var embeddedFiles embed.FS

// Download downloads the YouTube video/audio using yt-dlp.
func Download(url string, format int) error {
	// Extract embedded binaries to temp directory
	tempDir := os.TempDir()

	ytDlpPath, err := extractBinary("yt-dlp.exe", tempDir)
	if err != nil {
		return fmt.Errorf("failed to extract yt-dlp: %v", err)
	}

	if _, err := extractBinary("ffmpeg.exe", tempDir); err != nil {
		return fmt.Errorf("failed to extract ffmpeg: %v", err)
	}

	if _, err := extractBinary("ffprobe.exe", tempDir); err != nil {
		return fmt.Errorf("failed to extract ffprobe: %v", err)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("unable to get user home directory: %v", err)
	}
	downloadPath := filepath.Join(home, "Downloads")
	outputTemplate := filepath.Join(downloadPath, "%(title).80s.%(ext)s")

	args := []string{
		url,
		"--ffmpeg-location", tempDir,
		"-o", outputTemplate,
	}

	// Format: 1=mp3, 2=480p, 3=720p
	switch format {
	case 1:
		// Extract audio only as mp3, keep original video file? (add -k if you want)
		args = append(args, "-x", "--audio-format", "mp3")
	case 2:
		// Best video up to 480p + best audio, merge as mp4
		args = append(args, "-f", "bestvideo[height<=480]+bestaudio/best[height<=480]", "--merge-output-format", "mp4")
	case 3:
		// Best video up to 720p + best audio, merge as mp4
		args = append(args, "-f", "bestvideo[height<=720]+bestaudio/best[height<=720]", "--merge-output-format", "mp4")
	case 4:
		// Best video up to 1080p + best audio, merge as mp4
		args = append(args, "-f", "bestvideo[height<=1080]+bestaudio/best[height<=1080]", "--merge-output-format", "mp4")
	default:
		return fmt.Errorf("unsupported format selection")
	}

	cmd := exec.Command(ytDlpPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("Downloading... this may take some time.")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("download failed: %v", err)
	}

	return nil
}

func extractBinary(name, destDir string) (string, error) {
	data, err := embeddedFiles.ReadFile(name)
	if err != nil {
		return "", err
	}

	outPath := filepath.Join(destDir, name)
	if err := os.WriteFile(outPath, data, 0755); err != nil {
		return "", err
	}

	return outPath, nil
}
