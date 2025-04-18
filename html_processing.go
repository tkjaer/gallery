package main

import (
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"text/template"
)

// processHTMLFile is called when an HTML file is found that needs to be processed.
// It will index the directory content, and generate a new index file in the
// output directory.
func processHTMLFile(htmlTasks <-chan Dir, wg *sync.WaitGroup, done <-chan struct{}) {
	slog.Debug("Starting processHTML goroutine")

	tpl, err := template.ParseGlob("templates/default/*.go.html")
	if err != nil {
		slog.Error("Failed to parse template", "error", err)
		return
	}
	slog.Debug("Template parsed", "template", tpl)

	for {
		select {
		case htmlTask := <-htmlTasks:
			defer wg.Done()
			slog.Debug("Received HTML task", "htmlTask", htmlTask)

			navigation := []NavigationElement{}
			images := []Image{}
			folders := []string{}
			fileIndex := 0

			imagePath := strings.TrimPrefix(htmlTask.Path, config.Originals)
			imagePath = strings.TrimPrefix(imagePath, "/")
			outputDir := filepath.Join(config.Output, imagePath)
			outputFile := filepath.Join(outputDir, "index.html")

			for _, image := range htmlTask.Files {
				fileIndex += 1
				images = append(images, Image{
					Description: image.Name,
					File:        image.Name,
					Path:        imagePath,
					Index:       fileIndex,
				})
				slog.Debug("Image added", "image", image.Name, "path", imagePath)
			}

			navigationParts := strings.Split(imagePath, "/")
			for i := range navigationParts {
				slog.Debug("Processing navigation part", "navigationPart", navigationParts[i])
				navigation = append(navigation, NavigationElement{
					Path: strings.Join(navigationParts[:i+1], "/"),
					Name: navigationParts[i],
				})
				slog.Debug("Directory added", "path", navigationParts[i])
			}

			for _, subDir := range htmlTask.SubDirs {
				slog.Debug("Processing subdirectory", "subDir", subDir.Name)
				folders = append(folders, subDir.Name)
				slog.Debug("Subdirectory added", "subDir", subDir.Name)
			}

			// Sort images based on the specified order
			switch config.ImageOrder {
			case "alphabetical":
				slog.Debug("Sorting images alphabetically")
				// Images are already in alphabetical order by default
			case "new":
				slog.Debug("Sorting images by newest first")
				// Sort images by ModTime in descending order
				sort.SliceStable(images, func(i, j int) bool {
					return htmlTask.Files[images[i].File].ModTime.After(htmlTask.Files[images[j].File].ModTime)
				})
			case "old":
				slog.Debug("Sorting images by oldest first")
				// Sort images by ModTime in ascending order
				sort.SliceStable(images, func(i, j int) bool {
					return htmlTask.Files[images[i].File].ModTime.Before(htmlTask.Files[images[j].File].ModTime)
				})
			}

			g := Gallery{
				Name:        config.Name,
				Copyright:   config.Copyright,
				Folders:     folders,
				Navigation:  navigation,
				Images:      images,
				Year:        year,
				GalleryPath: config.GalleryPath,
			}
			slog.Debug("Gallery object created", "gallery", g)

			err := os.MkdirAll(filepath.Join(outputDir), 0755)
			if err != nil {
				slog.Error("Failed to create output directory", "error", err)
				os.Exit(1)
			}
			slog.Debug("Output directory created", "outputDir", outputDir)

			f, err := os.Create(outputFile)
			if err != nil {
				slog.Error("Failed to create output file", "error", err)
				os.Exit(1)
			}
			defer f.Close()
			slog.Debug("Output file created", "outputFile", outputFile)

			err = tpl.ExecuteTemplate(f, "index.go.html", g)
			if err != nil {
				slog.Error("Failed to execute template", "error", err)
				os.Exit(1)
			}
			slog.Debug("Template executed", "outputFile", outputFile)

		case <-done:
			slog.Debug("Received done signal")
			return
		}
	}
}
