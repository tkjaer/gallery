package main

import (
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcess(t *testing.T) {
	// Set up temporary directories for testing
	tempDir := t.TempDir()
	config.Output = filepath.Join(tempDir, "output")
	config.Originals = filepath.Join(tempDir, "originals")
	config.ThumbSize = 100
	config.FullSize = 800
	config.JPEGQuality = 90
	config.CopyOriginals = false

	// Create the originals directory
	err := os.MkdirAll(config.Originals, 0755)
	assert.NoError(t, err)

	// Create mock image files in the originals directory
	imagePaths := []string{
		filepath.Join(config.Originals, "image1.jpg"),
		filepath.Join(config.Originals, "image2.jpg"),
	}
	for _, imagePath := range imagePaths {
		img := image.NewRGBA(image.Rect(0, 0, 200, 100))
		file, err := os.Create(imagePath)
		assert.NoError(t, err)
		defer file.Close()

		err = jpeg.Encode(file, img, &jpeg.Options{Quality: 90})
		assert.NoError(t, err)
	}

	// Create a mock subdirectory in the originals directory
	subDir := filepath.Join(config.Originals, "subdir")
	err = os.MkdirAll(subDir, 0755)
	assert.NoError(t, err)

	// Create a mock image file in the subdirectory
	subImagePath := filepath.Join(subDir, "subimage.jpg")
	file, err := os.Create(subImagePath)
	assert.NoError(t, err)
	defer file.Close()

	img := image.NewRGBA(image.Rect(0, 0, 200, 100))
	err = jpeg.Encode(file, img, &jpeg.Options{Quality: 90})
	assert.NoError(t, err)

	// Run the process function
	err = process()
	assert.NoError(t, err)

	// Verify that the output directory was created
	_, err = os.Stat(config.Output)
	assert.NoError(t, err)

	// Verify that the thumbnails and full images were created
	for _, imagePath := range imagePaths {
		imageName := filepath.Base(imagePath)
		thumbPath := filepath.Join(config.Output, "thumb_"+imageName)
		fullPath := filepath.Join(config.Output, "full_"+imageName)

		_, err = os.Stat(thumbPath)
		assert.NoError(t, err)

		_, err = os.Stat(fullPath)
		assert.NoError(t, err)
	}

	// Verify that the subdirectory was processed
	subOutputDir := filepath.Join(config.Output, "subdir")
	_, err = os.Stat(subOutputDir)
	assert.NoError(t, err)

	subThumbPath := filepath.Join(subOutputDir, "thumb_subimage.jpg")
	subFullPath := filepath.Join(subOutputDir, "full_subimage.jpg")

	_, err = os.Stat(subThumbPath)
	assert.NoError(t, err)

	_, err = os.Stat(subFullPath)
	assert.NoError(t, err)
}

func TestProcessWithCopyOriginals(t *testing.T) {
	// Set up temporary directories for testing
	tempDir := t.TempDir()
	config.Output = filepath.Join(tempDir, "output")
	config.Originals = filepath.Join(tempDir, "originals")
	config.CopyOriginals = true

	// Create the originals directory
	err := os.MkdirAll(config.Originals, 0755)
	assert.NoError(t, err)

	// Create mock image files in the originals directory
	imagePaths := []string{
		filepath.Join(config.Originals, "image1.jpg"),
		filepath.Join(config.Originals, "image2.jpg"),
	}
	for _, imagePath := range imagePaths {
		img := image.NewRGBA(image.Rect(0, 0, 200, 100))
		file, err := os.Create(imagePath)
		assert.NoError(t, err)
		defer file.Close()

		err = jpeg.Encode(file, img, &jpeg.Options{Quality: 90})
		assert.NoError(t, err)
	}

	// Run the process function
	err = process()
	assert.NoError(t, err)

	// Verify that the original images were copied to the output directory
	for _, imagePath := range imagePaths {
		imageName := filepath.Base(imagePath)
		fullPath := filepath.Join(config.Output, "full_"+imageName)

		_, err = os.Stat(fullPath)
		assert.NoError(t, err)
	}
}
