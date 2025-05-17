package main

import (
	"github.com/sanjeevnode/go-video-downloader/internal/config"
	"github.com/sanjeevnode/go-video-downloader/internal/menu"
)

func main() {
	config.LoadEnv() // loads .env
	menu.ShowMainMenu()
}
