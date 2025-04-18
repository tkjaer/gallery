package main

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
)

// processImage is called when an image is found that needs to be processed.
// It will resize the image, and copy it to the output directory.
func processImage(imageTasks <-chan string, wg *sync.WaitGroup, done <-chan struct{}) {
	slog.Debug("Starting processImage goroutine")
	var file string
	for {
		select {
		case file = <-imageTasks:
			defer wg.Done()
			if file == "" {
				slog.Debug("Received empty file path, skipping")
				continue
			}
			slog.Debug("Received image task", "file", file)

			imgName := filepath.Base(file)
			outputDir := filepath.Join(config.Output, filepath.Dir(strings.TrimPrefix(file, config.Originals)))

			img, err := imgio.Open(file)
			if err != nil {
				slog.Error("Failed to open image", "file", file, "error", err)
				os.Exit(1)
			}
			slog.Debug("Image opened", "file", file)
			// calculate height, depending on the aspect ratio and the config.ThumbSize
			width := img.Bounds().Max.X
			height := img.Bounds().Max.Y
			slog.Debug("Image dimensions", "width", width, "height", height)
			aspectRatio := float64(width) / float64(height)
			thumbWidth := config.ThumbSize
			thumbHeight := int(float64(thumbWidth) / aspectRatio)
			slog.Debug("Aspect ratio calculated", "aspectRatio", aspectRatio, "thumbWidth", thumbWidth, "thumbHeight", thumbHeight)

			err = os.MkdirAll(outputDir, 0755)
			if err != nil {
				slog.Error("Failed to create output directory", "error", err)
				os.Exit(1)
			}
			slog.Debug("Output directory created", "outputDir", outputDir)

			// Generate thumbnail
			thumb := transform.Resize(img, config.ThumbSize, thumbHeight, transform.Lanczos)
			slog.Debug("Thumbnail resized", "thumbSize", config.ThumbSize)
			if err := imgio.Save(filepath.Join(outputDir, "thumb_"+imgName), thumb, imgio.JPEGEncoder(config.JPEGQuality)); err != nil {
				slog.Error("Failed to save thumbnail", "error", err)
				os.Exit(1)
			}
			slog.Debug("Thumbnail saved", "thumbFile", filepath.Join(outputDir, "thumb_"+imgName))

			// Copy original or generate full image
			if config.CopyOriginals {
				slog.Debug("Copying original file", "file", file)
				err := copyFile(file, filepath.Join(outputDir, "full_"+imgName))
				if err != nil {
					slog.Error("Failed to copy original file", "error", err)
					os.Exit(1)
				}
				slog.Debug("Original file copied", "file", file)

			} else {
				full := transform.Resize(img, config.FullSize, config.FullSize, transform.Linear)
				slog.Debug("Full image resized", "fullSize", config.FullSize)
				if err := imgio.Save(filepath.Join(outputDir, "full_"+imgName), full, imgio.JPEGEncoder(config.JPEGQuality)); err != nil {
					slog.Error("Failed to save full image", "error", err)
					os.Exit(1)
				}
				slog.Debug("Full image saved", "fullFile", filepath.Join(outputDir, "full_"+imgName))
			}

		case <-done:
			slog.Debug("Received done signal")
			return
		}
	}
}
