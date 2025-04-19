package main

import (
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"sort"
	"sync"
	"text/template"
	"time"
)

func processRSSFeed(rssTasks <-chan RSSItem, wg *sync.WaitGroup, done <-chan struct{}) {
	slog.Debug("Starting processRSSFeed goroutine")
	defer wg.Done()

	tpl, err := template.ParseGlob(filepath.Join("templates", config.Template, "rss.go.xml"))
	if err != nil {
		slog.Error("Failed to parse template", "error", err)
		return
	}
	slog.Debug("Template parsed", "template", tpl)

	RSSFeed := RSSFeed{
		Title:         config.Name,
		Description:   "Latest images from " + config.Name,
		Link:          config.GalleryURL + config.GalleryPath,
		Language:      "en-us",
		Copyright:     config.Copyright,
		AtomLink:      config.GalleryURL + path.Join(config.GalleryPath, "rss.xml"),
		LastBuildDate: time.Now().Format(time.RFC1123Z),
		Items:         []RSSItem{},
	}
	rssFile := filepath.Join(config.Output, "rss.xml")
	updateRSSFeed := false

	for {
		select {
		case item := <-rssTasks:
			// Process the RSS item and add it to the RSS feed
			if item == (RSSItem{}) {
				slog.Debug("Received empty RSS item, skipping")
				continue
			}
			RSSFeed.Items = append(RSSFeed.Items, item)
			slog.Debug("RSS item added", "item", item)
		case <-done:
			// Sort the RSS items by PubDate
			if !config.RSSFeed {
				slog.Debug("RSS feed generation is disabled, skipping")
				return
			}
			slog.Debug("Sorting RSS items by PubDate", "itemCount", len(RSSFeed.Items))
			sort.Slice(RSSFeed.Items, func(i, j int) bool {
				pubDateI, errI := time.Parse(time.RFC1123Z, RSSFeed.Items[i].PubDate)
				if errI != nil {
					slog.Error("Failed to parse PubDate for item i", "error", errI)
					return false
				}
				pubDateJ, errJ := time.Parse(time.RFC1123Z, RSSFeed.Items[j].PubDate)
				if errJ != nil {
					slog.Error("Failed to parse PubDate for item j", "error", errJ)
					return false
				}
				return pubDateI.After(pubDateJ)
			})
			// Check if the RSS feed file exists and if it has been modified since the newest item PubDate
			if len(RSSFeed.Items) > 0 {
				slog.Debug("Checking if RSS feed file needs to be updated")
				rssFileInfo, err := os.Stat(rssFile)
				if err == nil {
					// File exists, check modification time
					rssFileModTime := rssFileInfo.ModTime()
					pubDate, err := time.Parse(time.RFC1123Z, RSSFeed.Items[0].PubDate)
					if err != nil {
						slog.Error("Failed to parse PubDate for comparison", "error", err)
						return
					}
					if rssFileModTime.After(pubDate) {
						slog.Debug("RSS feed file is up to date, skipping write")
						return
					} else {
						slog.Debug("RSS feed file is outdated, updating")
						updateRSSFeed = true
					}
				}
				// If the file does not exist, we need to create it
				if os.IsNotExist(err) {
					slog.Debug("RSS feed file does not exist, creating")
					updateRSSFeed = true
				} else if err != nil {
					slog.Error("Failed to stat RSS feed file", "error", err)
					return
				}
			}
			if updateRSSFeed {
				slog.Debug("Updating RSS feed file")
				// Only include the 100 newest items in the RSS feed
				if len(RSSFeed.Items) > 100 {
					slog.Debug("Trimming RSS feed items to 100 items")
					RSSFeed.Items = RSSFeed.Items[:100]
				}
				// Render the RSS feed to a file
				f, err := os.Create(rssFile)
				if err != nil {
					slog.Error("Failed to create RSS feed file", "error", err)
					return
				}
				defer f.Close()
				slog.Debug("RSS feed file created", "rssFile", rssFile)
				err = tpl.ExecuteTemplate(f, "rss.go.xml", RSSFeed)
				if err != nil {
					slog.Error("Failed to execute template for RSS feed", "error", err)
					return
				}
				slog.Debug("RSS feed file written", "rssFile", rssFile)
			}

			// Render the RSS feed to a file and return
			// If the output RSS feed exists, and has been modified since the newest item PubDate, skip writing and return
			slog.Debug("Received done signal, exiting processRSSFeed goroutine")
			return
		}
	}
}
