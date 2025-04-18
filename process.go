package main

import (
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// process walks the original directory, processes images, and generates HTML files for each directory.
func process() error {
	slog.Debug("Processing content")

	err := checkOrCreateOutputDir()
	if err != nil {
		slog.Error("Failed to check or create output directory", "error", err)
		os.Exit(1)
	}

	done := make(chan struct{})
	imageTasks := make(chan string)
	htmlTasks := make(chan Dir)

	defer close(imageTasks)
	defer close(htmlTasks)

	wg := &sync.WaitGroup{}
	slog.Debug("Created wait group", "waitGroup", wg)

	// numRoutines := runtime.NumCPU()
	numRoutines := 1

	// Start the image processing goroutines
	for range numRoutines {
		slog.Debug("Starting image processing goroutines", "numRoutines", numRoutines)
		wg.Add(1)
		go processImage(imageTasks, wg, done)
	}

	// Start the HTML processing goroutines
	for range numRoutines {
		slog.Debug("Starting HTML processing goroutines", "numRoutines", numRoutines)
		wg.Add(1)
		go processHTMLFile(htmlTasks, wg, done)
	}

	galleryContent := DirMap{}
	// Walk the original directory and send image tasks to the channel
	err = filepath.WalkDir(config.Originals, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		slog.Debug("Processing path", "path", path)

		name := d.Name()
		fileInfo, err := d.Info()
		if err != nil {
			return err
		}
		modTime := fileInfo.ModTime()

		parentDir := filepath.Dir(path)
		outputDir := filepath.Join(config.Output, strings.TrimPrefix(parentDir, config.Originals))

		needsUpdate := false

		if d.IsDir() {
			slog.Debug("Processing directory", "path", path, "name", name)
			outputIndex := filepath.Join(outputDir, "index.html")

			// Check if the output directory exists and if it's newer than the original directory
			// If it doesn't exist or is older, we need to update the HTML file for it
			outputIndexInfo, err := os.Stat(outputIndex)
			if err != nil {
				if os.IsNotExist(err) {
					slog.Debug("Output index does not exist", "outputIndex", outputIndex)
					needsUpdate = true
				}
			} else {
				if modTime.After(outputIndexInfo.ModTime()) {
					slog.Debug("Original directory is newer", "originalDir", path, "outputIndex", outputIndex)
					needsUpdate = true
				} else {
					slog.Debug("Output index is newer", "originalDir", path, "outputIndex", outputIndex)
				}
			}

			galleryContent.AddDir(path, name, needsUpdate)

			// Add the directory to the parent directory's subdirectories
			if _, ok := galleryContent[parentDir]; ok {
				slog.Debug("Adding subdirectory", "path", path, "name", name)
				galleryContent[parentDir].SubDirs[path] = SubDir{
					Name: name,
				}
			}
		} else {
			slog.Debug("Processing file", "path", path, "name", name)
			if strings.HasSuffix(name, ".jpg") || strings.HasSuffix(name, ".jpeg") {
				for _, size := range []string{"thumb", "full"} {
					outputFile := filepath.Join(outputDir, size+"_"+name)
					outputFileInfo, err := os.Stat(outputFile)
					if err != nil {
						if os.IsNotExist(err) {
							slog.Debug("Output file does not exist", "outputFile", outputFile)
							needsUpdate = true
						} else {
							slog.Error("Failed to stat output file", "error", err)
							os.Exit(1)
						}
					} else {
						if modTime.After(outputFileInfo.ModTime()) {
							slog.Debug("Original file is newer", "originalFile", path, "outputFile", outputFile)
							needsUpdate = true
						} else {
							slog.Debug("Output file is newer", "originalFile", path, "outputFile", outputFile)
						}
					}
				}
				if needsUpdate {
					imageTasks <- path
				}
				slog.Debug("Adding file to directory index", "path", path, "name", name)
				galleryContent[parentDir].Files[path] = File{
					Name:    name,
					ModTime: modTime,
				}
			} else {
				slog.Debug("Ignoring non-jpg file", "path", path)
			}
		}
		return nil
	})
	if err != nil {
		slog.Error("Failed to walk original directory", "error", err)
		os.Exit(1)
	}

	for _, dir := range galleryContent {
		if dir.NeedsUpdate {
			slog.Debug("Processing directory", "dir", dir)
			if len(dir.Files) > 0 || len(dir.SubDirs) > 0 {
				slog.Debug("Adding directory to HTML tasks", "dir", dir)
				htmlTasks <- dir
			}
		}
	}

	close(done)
	wg.Wait()

	err = updateTemplateFiles()
	if err != nil {
		return err
	}

	slog.Debug("Processing completed")
	return nil
}
