package main

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestProcessHTMLFile(t *testing.T) {
	// Set up temporary directories for testing
	tempDir := t.TempDir()
	config.Output = filepath.Join(tempDir, "output")
	config.Originals = filepath.Join(tempDir, "originals")
	config.Name = "Test Gallery"
	config.Copyright = "© 2025 Test"
	config.GalleryPath = "/gallery"
	config.ImageOrder = "new"

	// Create a mock template file
	templateDir := filepath.Join(tempDir, "templates", "default")
	err := os.MkdirAll(templateDir, 0755)
	assert.NoError(t, err)

	templateFile := filepath.Join(templateDir, "index.go.html")
	err = os.WriteFile(templateFile, []byte(`
        <html>
        <head><title>{{.Name}}</title></head>
        <body>
            <h1>{{.Name}}</h1>
            <ul>
                {{range .Images}}
                    <li>{{.}}</li>
                {{end}}
            </ul>
            <p>{{.Copyright}}</p>
        </body>
        </html>
    `), 0644)
	assert.NoError(t, err)

	// Set up channels and WaitGroup
	htmlTasks := make(chan Dir)
	done := make(chan struct{})
	var wg sync.WaitGroup

	// Start the processHTMLFile function in a goroutine
	wg.Add(1)
	go processHTMLFile(htmlTasks, &wg, done)

	// Create a mock HTML task
	htmlTask := Dir{
		Path: "/test",
		Files: map[string]File{
			"image1.jpg": {Name: "image1.jpg"},
			"image2.jpg": {Name: "image2.jpg"},
		},
		SubDirs: map[string]SubDir{
			"subdir1": {Name: "subdir1"},
		},
	}

	// Send the task to the channel
	htmlTasks <- htmlTask

	close(done)
	// Wait for the goroutine to finish
	wg.Wait()

	// Verify the output directory was created
	outputDir := filepath.Join(config.Output, "test")
	_, err = os.Stat(outputDir)
	assert.NoError(t, err)

	// Verify the output file was created
	outputFile := filepath.Join(outputDir, "index.html")
	_, err = os.Stat(outputFile)
	assert.NoError(t, err)

	// Verify the content of the output file
	content, err := os.ReadFile(outputFile)
	assert.NoError(t, err)
	assert.Contains(t, string(content), "Test Gallery")
	assert.Contains(t, string(content), "© 2025 Test")
	assert.Contains(t, string(content), "image1.jpg")
	assert.Contains(t, string(content), "image2.jpg")
}

func TestProcessHTMLFileWithNewestFirst(t *testing.T) {
	// Set up temporary directories for testing
	tempDir := t.TempDir()
	config.Output = filepath.Join(tempDir, "output")
	config.Originals = filepath.Join(tempDir, "originals")
	config.Name = "Test Gallery"
	config.Copyright = "© 2025 Test"
	config.GalleryPath = "/gallery"
	config.ImageOrder = "new"

	// Create a mock template file
	templateDir := filepath.Join(tempDir, "templates", "default")
	err := os.MkdirAll(templateDir, 0755)
	assert.NoError(t, err)

	templateFile := filepath.Join(templateDir, "index.go.html")
	err = os.WriteFile(templateFile, []byte(`
        <html>
        <head><title>{{.Name}}</title></head>
        <body>
            <h1>{{.Name}}</h1>
            <ul>
                {{range .Images}}
                    <li>{{.}}</li>
                {{end}}
            </ul>
            <p>{{.Copyright}}</p>
        </body>
        </html>
    `), 0644)
	assert.NoError(t, err)

	// Set up channels and WaitGroup
	htmlTasks := make(chan Dir)
	done := make(chan struct{})
	var wg sync.WaitGroup

	// Start the processHTMLFile function in a goroutine
	wg.Add(1)
	go processHTMLFile(htmlTasks, &wg, done)

	// Create a mock HTML task
	htmlTask := Dir{
		Path: "/test",
		Files: map[string]File{
			"image1.jpg": {
				Name:    "image1.jpg",
				ModTime: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			"image2.jpg": {
				Name:    "image2.jpg",
				ModTime: time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		SubDirs: map[string]SubDir{
			"subdir1": {Name: "subdir1"},
		},
	}

	// Send the task to the channel
	htmlTasks <- htmlTask

	close(done)
	// Wait for the goroutine to finish
	wg.Wait()

	// Verify the output directory was created
	outputDir := filepath.Join(config.Output, "test")
	_, err = os.Stat(outputDir)
	assert.NoError(t, err)

	// Verify the output file was created
	outputFile := filepath.Join(outputDir, "index.html")
	_, err = os.Stat(outputFile)
	assert.NoError(t, err)

	// Verify the content of the output file
	content, err := os.ReadFile(outputFile)
	assert.NoError(t, err)
	assert.Contains(t, string(content), "Test Gallery")
	assert.Contains(t, string(content), "© 2025 Test")
	assert.Contains(t, string(content), "image1.jpg")
	assert.Contains(t, string(content), "image2.jpg")

	// Verify the order of images (newest first)
	assert.True(t, strings.Index(string(content), "image2.jpg") < strings.Index(string(content), "image1.jpg"))
}
