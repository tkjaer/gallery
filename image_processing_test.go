package main

import (
	"image"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/stretchr/testify/assert"
)

func TestProcessImage(t *testing.T) {
	// Set up temporary directories for testing
	tempDir := t.TempDir()
	setupConfig(tempDir)

	// Create the originals directory
	err := os.MkdirAll(config.Originals, 0755)
	assert.NoError(t, err)

	// Create a mock image file in the originals directory
	originalImagePath := filepath.Join(config.Originals, "test.jpg")
	img := image.NewRGBA(image.Rect(0, 0, 200, 100))
	err = imgio.Save(originalImagePath, img, imgio.JPEGEncoder(90))
	assert.NoError(t, err)

	// Set up channels and WaitGroup
	imageTasks := make(chan string)
	rssTasks := make(chan RSSItem)
	rssDone := make(chan struct{})
	done := make(chan struct{})
	var wg sync.WaitGroup
	var rssWg sync.WaitGroup

	// Start the processRSSFeed function in a goroutine
	rssWg.Add(1)
	go processRSSFeed(rssTasks, &rssWg, rssDone)

	// Start the processImage function in a goroutine
	wg.Add(1)
	go processImage(imageTasks, rssTasks, &wg, done)

	// Add the image task to the channel
	imageTasks <- originalImagePath
	close(done)

	// Wait for the goroutine to finish
	wg.Wait()
	close(rssDone)
	rssWg.Wait()

	// Verify that the output directory was created
	outputDir := filepath.Join(config.Output, filepath.Dir(strings.TrimPrefix(originalImagePath, config.Originals)))
	_, err = os.Stat(outputDir)
	assert.NoError(t, err)

	// Verify that the thumbnail was created
	thumbPath := filepath.Join(outputDir, "thumb_test.jpg")
	_, err = os.Stat(thumbPath)
	assert.NoError(t, err)

	// Verify that the full image was created
	fullPath := filepath.Join(outputDir, "full_test.jpg")
	_, err = os.Stat(fullPath)
	assert.NoError(t, err)
}

func TestProcessImageWithCopyOriginals(t *testing.T) {
	// Set up temporary directories for testing
	tempDir := t.TempDir()
	config.Output = filepath.Join(tempDir, "output")
	config.Originals = filepath.Join(tempDir, "originals")
	config.ThumbSize = 100
	config.FullSize = 800
	config.JPEGQuality = 90
	config.CopyOriginals = true

	// Create the originals directory
	err := os.MkdirAll(config.Originals, 0755)
	assert.NoError(t, err)

	// Create a mock image file in the originals directory
	originalImagePath := filepath.Join(config.Originals, "test.jpg")
	img := image.NewRGBA(image.Rect(0, 0, 200, 100))
	err = imgio.Save(originalImagePath, img, imgio.JPEGEncoder(90))
	assert.NoError(t, err)

	// Set up channels and WaitGroup
	imageTasks := make(chan string, 1)
	rssTasks := make(chan RSSItem)
	done := make(chan struct{})
	rssDone := make(chan struct{})
	var wg sync.WaitGroup
	var rssWg sync.WaitGroup

	// Start the processRSSFeed function in a goroutine
	rssWg.Add(1)
	go processRSSFeed(rssTasks, &rssWg, rssDone)

	// Start the processImage function in a goroutine
	wg.Add(1)
	go processImage(imageTasks, rssTasks, &wg, done)

	// Add the image task to the channel
	imageTasks <- originalImagePath
	time.Sleep(10 * time.Millisecond) // Ensure the task is processed
	close(done)

	// Wait for the goroutine to finish
	wg.Wait()
	close(rssDone)
	rssWg.Wait()

	// Verify that the output directory was created
	outputDir := filepath.Join(config.Output, filepath.Dir(strings.TrimPrefix(originalImagePath, config.Originals)))
	_, err = os.Stat(outputDir)
	assert.NoError(t, err)

	// Verify that the thumbnail was created
	thumbPath := filepath.Join(outputDir, "thumb_test.jpg")
	_, err = os.Stat(thumbPath)
	assert.NoError(t, err)

	// Verify that the original file was copied
	fullPath := filepath.Join(outputDir, "full_test.jpg")
	_, err = os.Stat(fullPath)
	assert.NoError(t, err)
}

func setupConfig(tempDir string) {
	config.Output = filepath.Join(tempDir, "output")
	config.Originals = filepath.Join(tempDir, "originals")
	config.ThumbSize = 100
	config.FullSize = 800
	config.JPEGQuality = 90
	config.CopyOriginals = false
}
