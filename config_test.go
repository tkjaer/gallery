package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig_Defaults(t *testing.T) {
	// Reset the global config to defaults
	config = Config{
		Name:          "Photo Gallery",
		Originals:     "originals",
		Output:        "output",
		Template:      "default",
		ThumbSize:     200,
		FullSize:      2000,
		CopyOriginals: false,
		ImageOrder:    "old",
		JPEGQuality:   90,
		GalleryPath:   "/",
		RSSFeed:       true,
	}

	// Call LoadConfig with a non-existent file
	err := LoadConfig("nonexistent.yaml")
	assert.NoError(t, err)

	// Verify that the default configuration is used
	assert.Equal(t, "Photo Gallery", config.Name)
	assert.Equal(t, "originals", config.Originals)
	assert.Equal(t, "output", config.Output)
	assert.Equal(t, "default", config.Template)
	assert.Equal(t, 200, config.ThumbSize)
	assert.Equal(t, 2000, config.FullSize)
	assert.Equal(t, false, config.CopyOriginals)
	assert.Equal(t, "new", config.ImageOrder)
	assert.Equal(t, 90, config.JPEGQuality)
	assert.Equal(t, "/", config.GalleryPath)
	assert.Equal(t, false, config.RSSFeed)
}

func TestLoadConfig_ValidFile(t *testing.T) {
	// Create a temporary YAML configuration file
	tempFile, err := os.CreateTemp("", "config_*.yaml")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	configContent := `
name: Test Gallery
originals: test_originals
output: test_output
template: custom
thumbnail_size: 150
full_size: 1200
copy_originals: true
image_order: old
jpeg_quality: 80
gallery_path: /test
rss_feed: false
`
	_, err = tempFile.Write([]byte(configContent))
	assert.NoError(t, err)
	tempFile.Close()

	// Call LoadConfig with the temporary file
	err = LoadConfig(tempFile.Name())
	assert.NoError(t, err)

	// Verify that the configuration is loaded correctly
	assert.Equal(t, "Test Gallery", config.Name)
	assert.Equal(t, "test_originals", config.Originals)
	assert.Equal(t, "test_output", config.Output)
	assert.Equal(t, "custom", config.Template)
	assert.Equal(t, 150, config.ThumbSize)
	assert.Equal(t, 1200, config.FullSize)
	assert.Equal(t, true, config.CopyOriginals)
	assert.Equal(t, "old", config.ImageOrder)
	assert.Equal(t, 80, config.JPEGQuality)
	assert.Equal(t, "/test", config.GalleryPath)
	assert.Equal(t, false, config.RSSFeed)
}

func TestLoadConfig_MissingFile(t *testing.T) {
	// Call LoadConfig with a non-existent file
	err := LoadConfig("nonexistent.yaml")
	assert.NoError(t, err)

	// Verify that the default configuration is used
	assert.Equal(t, "Photo Gallery", config.Name)
	assert.Equal(t, "originals", config.Originals)
	assert.Equal(t, "output", config.Output)
}

func TestLoadConfig_InvalidYAML(t *testing.T) {
	// Create a temporary invalid YAML configuration file
	tempFile, err := os.CreateTemp("", "config_invalid_*.yaml")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	_, err = tempFile.Write([]byte("invalid_yaml: [this is not valid YAML"))
	assert.NoError(t, err)
	tempFile.Close()

	// Call LoadConfig with the invalid file
	err = LoadConfig(tempFile.Name())
	assert.Error(t, err)
}

func TestLoadConfig_InvalidConfig(t *testing.T) {
	// Create a temporary YAML configuration file with invalid values
	tempFile, err := os.CreateTemp("", "config_invalid_*.yaml")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	configContent := `
originals: same_path
output: same_path
`
	_, err = tempFile.Write([]byte(configContent))
	assert.NoError(t, err)
	tempFile.Close()

	// Call LoadConfig with the invalid file
	err = LoadConfig(tempFile.Name())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "the \"originals\" and \"output\" directories cannot be the same")
}
