package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestCopyFile(t *testing.T) {
	// Create temporary source and destination files
	srcFile, err := os.CreateTemp("", "src")
	if err != nil {
		t.Fatalf("Failed to create temp source file: %v", err)
	}
	defer os.Remove(srcFile.Name())

	destFile, err := os.CreateTemp("", "dest")
	if err != nil {
		t.Fatalf("Failed to create temp destination file: %v", err)
	}
	defer os.Remove(destFile.Name())

	// Write some content to the source file
	content := []byte("test content")
	if _, err := srcFile.Write(content); err != nil {
		t.Fatalf("Failed to write to source file: %v", err)
	}
	srcFile.Close()

	// Copy the file
	err = copyFile(srcFile.Name(), destFile.Name())
	if err != nil {
		t.Fatalf("copyFile failed: %v", err)
	}

	// Verify the content of the destination file
	destContent, err := ioutil.ReadFile(destFile.Name())
	if err != nil {
		t.Fatalf("Failed to read destination file: %v", err)
	}
	if string(destContent) != string(content) {
		t.Errorf("Content mismatch: got %s, want %s", destContent, content)
	}
}

func TestProcessHTML(t *testing.T) {
	// Setup temporary directories and mock config
	tempDir := t.TempDir()
	config.Output = filepath.Join(tempDir, "output")
	config.Originals = filepath.Join(tempDir, "originals")
	err := os.MkdirAll(config.Originals, 0755)
	if err != nil {
		t.Fatalf("Failed to create originals directory: %v", err)
	}

	// Create a mock HTML file in the originals directory
	htmlFile := filepath.Join(config.Originals, "index.html")
	err = ioutil.WriteFile(htmlFile, []byte("<html></html>"), 0644)
	if err != nil {
		t.Fatalf("Failed to create HTML file: %v", err)
	}

	// Mock original content
	originalContent := DirMap{
		config.Originals: {
			Name:    "originals",
			Files:   map[string]File{htmlFile: {Name: "index.html"}},
			SubDirs: map[string]SubDir{},
		},
	}

	// Call processHTML
	err = processHTML([]string{config.Originals}, &originalContent)
	if err != nil {
		t.Fatalf("processHTML failed: %v", err)
	}

	// Verify the output file was created
	outputFile := filepath.Join(config.Output, "index.html")
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Errorf("Output file not created: %s", outputFile)
	}
}

func TestProcessImages(t *testing.T) {
	// Setup temporary directories and mock config
	tempDir := t.TempDir()
	config.Output = filepath.Join(tempDir, "output")
	config.Originals = filepath.Join(tempDir, "originals")
	err := os.MkdirAll(config.Originals, 0755)
	if err != nil {
		t.Fatalf("Failed to create originals directory: %v", err)
	}

	// Create a mock image file in the originals directory
	imageFile := filepath.Join(config.Originals, "image.jpg")
	img := image.NewRGBA(image.Rect(0, 0, 100, 100)) // Create a 100x100 test image
	for x := 0; x < 100; x++ {
		for y := 0; y < 100; y++ {
			img.Set(x, y, color.RGBA{R: uint8(x), G: uint8(y), B: 255, A: 255}) // Set pixel colors
		}
	}
	file, err := os.Create(imageFile)
	if err != nil {
		t.Fatalf("Failed to create image file: %v", err)
	}
	defer file.Close()
	err = jpeg.Encode(file, img, nil) // Encode the image as JPEG
	if err != nil {
		t.Fatalf("Failed to encode image: %v", err)
	}

	// Call processImages
	err = processImages([]string{imageFile})
	if err != nil {
		t.Fatalf("processImages failed: %v", err)
	}

	// Verify the thumbnail and full-size images were created
	thumbFile := filepath.Join(config.Output, "thumb_image.jpg")
	fullFile := filepath.Join(config.Output, "full_image.jpg")
	if _, err := os.Stat(thumbFile); os.IsNotExist(err) {
		t.Errorf("Thumbnail not created: %s", thumbFile)
	}
	if _, err := os.Stat(fullFile); os.IsNotExist(err) {
		t.Errorf("Full-size image not created: %s", fullFile)
	}
}
