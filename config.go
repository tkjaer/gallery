package main

import (
	"fmt"
	"log"
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
	NewestFirst   bool   `yaml:"newest_first" default:"true"`
	JPEGQuality   int    `yaml:"jpeg_quality" default:"90"`
	GalleryPath   string `yaml:"gallery_path" default:"/"`
}

var config = Config{
	Name:          "Photo Gallery",
	Originals:     "originals",
	Output:        "output",
	Template:      "default",
	ThumbSize:     200,
	FullSize:      2000,
	CopyOriginals: false,
	NewestFirst:   true,
	JPEGQuality:   90,
	GalleryPath:   "/",
}

// LoadConfig loads the configuration from a file.
func LoadConfig(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("No config file found, using defaults")
		} else {
			return err
		}
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return err
	}
	if config.Originals == config.Output {
		return fmt.Errorf("the \"originals\" and \"output\" directories cannot be the same")
	}

	return nil
}
