package main

import (
	"fmt"
	"log/slog"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Name          string `yaml:"name" default:"Photo Gallery"`
	Copyright     string `yaml:"copyright" default:""`
	Originals     string `yaml:"originals" default:"originals"`
	Output        string `yaml:"output" default:"output"`
	Template      string `yaml:"template" default:"default"`
	ThumbSize     int    `yaml:"thumbnail_size" default:"200"`
	FullSize      int    `yaml:"full_size" default:"2000"`
	CopyOriginals bool   `yaml:"copy_originals" default:"false"`
	ImageOrder    string `yaml:"image_order" default:"new"`
	JPEGQuality   int    `yaml:"jpeg_quality" default:"90"`
	GalleryPath   string `yaml:"gallery_path" default:"/"`
	GalleryURL    string `yaml:"gallery_url" default:""`
	RSSFeed       bool   `yaml:"rss_feed" default:"false"`
}

var config Config

// LoadConfig loads the configuration from a file.
func LoadConfig(filename string) error {
	// Initialize config with default values
	slog.Debug("Loading config file", "filename", filename)
	config = Config{
		Name:          "Photo Gallery",
		Copyright:     "",
		Originals:     "originals",
		Output:        "output",
		Template:      "default",
		ThumbSize:     200,
		FullSize:      2000,
		CopyOriginals: false,
		ImageOrder:    "new",
		JPEGQuality:   90,
		GalleryPath:   "/",
		GalleryURL:    "",
		RSSFeed:       false,
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			slog.Info("No config file found, using defaults")
			return nil
		} else {
			return err
		}
	}

	slog.Debug("Config file found, parsing", "filename", filename)
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return err
	}

	// Validate that ImageOrder is one of the allowed values ("new", "old", "alphabetical")
	if config.ImageOrder != "new" && config.ImageOrder != "old" && config.ImageOrder != "alphabetical" {
		return fmt.Errorf("invalid image order: %s, must be one of: new, old, alphabetical", config.ImageOrder)
	}

	// GalleryURL is required if RSSFeed is enabled
	if config.RSSFeed && config.GalleryURL == "" {
		return fmt.Errorf("gallery_url is required when rss_feed is enabled")
	}
	// Check that the GalleryURL looks just somewhat like a URL
	// This is a very basic check, we might want to use a more robust URL validation
	isValidURL := func(url string) bool {
		return len(url) > 0 && (url[:7] == "http://" || url[:8] == "https://")
	}
	// Validate the GalleryURL if RSSFeed is enabled
	if config.RSSFeed && !isValidURL(config.GalleryURL) {
		return fmt.Errorf("invalid gallery_url: %s", config.GalleryURL)
	}

	slog.Debug("Config file parsed successfully", "config", config)
	if config.Originals == config.Output {
		return fmt.Errorf("the \"originals\" and \"output\" directories cannot be the same")
	}

	return nil
}
