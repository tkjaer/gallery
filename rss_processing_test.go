package main

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestProcessRSSFeed(t *testing.T) {
	// Set up temporary directories for testing
	tempDir := t.TempDir()
	config.Output = filepath.Join(tempDir, "output")
	config.Template = "default"
	config.Name = "Test Gallery"
	config.Copyright = "Test Author"
	config.GalleryURL = "https://example.com"
	config.GalleryPath = "/gallery/"
	config.RSSFeed = true

	// Create the output directory
	err := os.MkdirAll(config.Output, 0755)
	assert.NoError(t, err)

	// Create a mock RSS template file
	templateDir := filepath.Join(tempDir, "templates", config.Template)
	err = os.MkdirAll(templateDir, 0755)
	assert.NoError(t, err)

	templateFile := filepath.Join(templateDir, "rss.go.xml")
	err = os.WriteFile(templateFile, []byte(`
        <rss version="2.0">
            <channel>
                <title>{{.Title}}</title>
                <link>{{.Link}}</link>
                <description>{{.Description}}</description>
                {{range .Items}}
                <item>
                    <title>{{.Title}}</title>
                    <link>{{.Link}}</link>
                    <description>{{.Description}}</description>
                    <pubDate>{{.PubDate}}</pubDate>
                    <guid>{{.GUID}}</guid>
                </item>
                {{end}}
            </channel>
        </rss>
    `), 0644)
	assert.NoError(t, err)

	// Set up channels and WaitGroup
	rssTasks := make(chan RSSItem)
	done := make(chan struct{})
	var rssWg sync.WaitGroup

	// Start the processRSSFeed function in a goroutine
	rssWg.Add(1)
	go processRSSFeed(rssTasks, &rssWg, done)

	var wg sync.WaitGroup

	// Send mock RSS items to the channel
	wg.Add(1)
	go func() {
		defer wg.Done()
		rssTasks <- RSSItem{
			Title:       "Image 1",
			Description: "Description for Image 1",
			Link:        "https://example.com/gallery/image1.jpg",
			PubDate:     time.Now().Add(-1 * time.Hour).Format(time.RFC1123Z),
			GUID:        "image1",
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		rssTasks <- RSSItem{
			Title:       "Image 2",
			Description: "Description for Image 2",
			Link:        "https://example.com/gallery/image2.jpg",
			PubDate:     time.Now().Add(-2 * time.Hour).Format(time.RFC1123Z),
			GUID:        "image2",
		}
	}()

	// Wait for the goroutines to finish sending items
	wg.Wait()

	// Close the rssTasks done channel
	close(done)
	// Wait for the goroutine to finish
	rssWg.Wait()

	// Verify that the RSS feed file was created
	rssFile := filepath.Join(config.Output, "rss.xml")
	_, err = os.Stat(rssFile)
	assert.NoError(t, err)

	// Verify the content of the RSS feed file
	content, err := os.ReadFile(rssFile)
	assert.NoError(t, err)
	assert.Contains(t, string(content), "<title>Test Gallery</title>")
	assert.Contains(t, string(content), "<link>https://example.com/gallery/</link>")
	assert.Contains(t, string(content), "<description>Latest images from Test Gallery</description>")
	assert.Contains(t, string(content), "<title>Image 1</title>")
	assert.Contains(t, string(content), "<title>Image 2</title>")
}

func TestProcessRSSFeed_Disabled(t *testing.T) {
	// Set up temporary directories for testing
	tempDir := t.TempDir()
	config.Output = filepath.Join(tempDir, "output")
	config.Template = "default"
	config.RSSFeed = false // RSS feed generation is disabled

	// Create the output directory
	err := os.MkdirAll(config.Output, 0755)
	assert.NoError(t, err)

	// Set up channels and WaitGroup
	rssTasks := make(chan RSSItem, 10)
	done := make(chan struct{})
	var wg sync.WaitGroup

	// Start the processRSSFeed function in a goroutine
	wg.Add(1)
	go processRSSFeed(rssTasks, &wg, done)

	// Send mock RSS items to the channel
	rssTasks <- RSSItem{
		Title:       "Image 1",
		Description: "Description for Image 1",
		Link:        "https://example.com/gallery/image1.jpg",
		PubDate:     time.Now().Add(-1 * time.Hour).Format(time.RFC1123Z),
		GUID:        "image1",
	}

	// Close the rssTasks channel and signal done
	close(rssTasks)
	close(done)

	// Wait for the goroutine to finish
	wg.Wait()

	// Verify that the RSS feed file was not created
	rssFile := filepath.Join(config.Output, "rss.xml")
	_, err = os.Stat(rssFile)
	assert.Error(t, err) // File should not exist
	assert.True(t, os.IsNotExist(err))
}
