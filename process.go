package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"text/template"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
)

// processHTML finds the HTML files that need to be processed, and launches parallel
// workers to process them via the processHTMLFile function.
func processHTML(indexUpdates []string, originalContent *DirMap) error {
	if len(indexUpdates) == 0 {
		return nil
	}
	log.Println("Processing HTML")

	for _, index := range indexUpdates {
		err := processHTMLFile(index, originalContent)
		if err != nil {
			return err
		}
	}
	return nil
}

// processHTMLFile is called when an HTML file is found that needs to be processed.
// It will index the directory content, and generate a new index file in the
// output directory.
func processHTMLFile(index string, originalContent *DirMap) error {
	log.Printf("Processing index %s", index)

	dirContent := (*originalContent)[index]
	imagePath := strings.TrimPrefix(index, config.Originals)
	imagePath = strings.TrimPrefix(imagePath, "/")
	outputDir := filepath.Join(config.Output, imagePath)
	outputFile := filepath.Join(outputDir, "index.html")

	directories := []Directory{}
	images := []Image{}
	folders := []string{}

	fileIndex := 0
	for _, image := range dirContent.Files {
		fileIndex += 1
		images = append(images, Image{
			Description: image.Name,
			File:        image.Name,
			Path:        imagePath,
			Index:       fileIndex,
		})
	}
	pathParts := strings.Split(imagePath, "/")
	for i := 0; i < len(pathParts); i++ {
		directories = append(directories, Directory{
			Path: strings.Join(pathParts[:i+1], "/"),
			Name: pathParts[i],
		})
	}

	for _, subDir := range dirContent.SubDirs {
		folders = append(folders, subDir.Name)
	}

	if config.NewestFirst {
		imagesReversed := make([]Image, len(images))
		for i, j := 0, len(images)-1; i < len(images); i, j = i+1, j-1 {
			imagesReversed[i] = images[j]
		}
		images = imagesReversed
	}

	g := Gallery{
		Name:        config.Name,
		Copyright:   config.Copyright,
		Folders:     folders,
		Directories: directories,
		Images:      images,
		Year:        year,
		GalleryPath: config.GalleryPath,
	}
	err := os.MkdirAll(filepath.Join(outputDir), 0755)
	if err != nil {
		return err
	}

	f, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer f.Close()

	tpl, err := template.ParseGlob("templates/default/*.go.html")
	if err != nil {
		return err
	}
	err = tpl.ExecuteTemplate(f, "index.go.html", g)
	if err != nil {
		return err
	}
	return nil
}

// processImages finds the images that need to be processed, and launches parallel
// workers to process them via the processImage function.
func processImages(fileUpdates []string) error {
	if len(fileUpdates) == 0 {
		return nil
	}
	log.Println("Processing images")

	runtime.GOMAXPROCS(runtime.NumCPU())
	wg := &sync.WaitGroup{}
	wg.Add(len(fileUpdates))
	for _, file := range fileUpdates {
		go func(file string) {
			defer wg.Done()
			err := processImage(file, config)
			if err != nil {
				log.Fatalf("Failed to process image: %v", err)
			}
		}(file)
	}
	wg.Wait()

	return nil
}

// copyFile copies a file from source to destination.
func copyFile(source, destination string) error {
	sourceFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err == nil {
		return err
	}
	return nil
}

// processImage is called when an image is found that needs to be processed.
// It will resize the image, and copy it to the output directory.
func processImage(file string, config Config) error {
	// log.Printf("Processing image %s", file)
	imgName := filepath.Base(file)
	outputDir := filepath.Join(config.Output, filepath.Dir(strings.TrimPrefix(file, config.Originals)))

	img, err := imgio.Open(file)
	if err != nil {
		return err
	}

	// calculate height, depending on the aspect ratio and the config.ThumbSize
	width := img.Bounds().Max.X
	height := img.Bounds().Max.Y
	aspectRatio := float64(width) / float64(height)
	thumbWidth := config.ThumbSize
	thumbHeight := int(float64(thumbWidth) / aspectRatio)

	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		return err
	}

	thumb := transform.Resize(img, config.ThumbSize, thumbHeight, transform.Lanczos)
	if err := imgio.Save(filepath.Join(outputDir, "thumb_"+imgName), thumb, imgio.JPEGEncoder(config.JPEGQuality)); err != nil {
		return err
	}

	if config.CopyOriginals {
		err := copyFile(file, filepath.Join(outputDir, "full_"+imgName))
		if err != nil {
			return err
		}
	} else {
		full := transform.Resize(img, config.FullSize, config.FullSize, transform.Linear)
		if err := imgio.Save(filepath.Join(outputDir, "full_"+imgName), full, imgio.JPEGEncoder(config.JPEGQuality)); err != nil {
			return err
		}
	}

	return nil
}

// getUpdates returns a list of index and file updates that are needed, based on the
// original and output content.
func getUpdates(originalContent DirMap, outputContent outputMap) ([]string, []string, error) {
	indexUpdates := []string{}
	fileUpdates := []string{}
	for origPath, origDetails := range originalContent {
		indexUpdateNeeded := true
		if outPath, ok := outputContent[config.Output+strings.TrimPrefix(origPath, config.Originals)]; ok {
			if outPath.After(origDetails.ModTime) {
				indexUpdateNeeded = false
			}
		}
		for origFilePath, origFileDetails := range origDetails.Files {
			if !strings.HasSuffix(origFilePath, ".jpg") && !strings.HasSuffix(origFilePath, ".jpeg") {
				continue
			}
			// Also update index if any image updates are required
			imageUpdateNeeded := false
			for _, outputSize := range []string{"thumb", "full"} {
				o := filepath.Join(
					config.Output,
					filepath.Dir(strings.TrimPrefix(origFilePath, config.Originals)),
					outputSize+"_"+origFileDetails.Name,
				)
				if outFilePath, ok := outputContent[o]; ok {
					if outFilePath.Before(origFileDetails.ModTime) {
						imageUpdateNeeded = true
						indexUpdateNeeded = true
					}
				} else {
					imageUpdateNeeded = true
					indexUpdateNeeded = true
				}
			}
			if imageUpdateNeeded {
				fileUpdates = append(fileUpdates, origFilePath)
			}
		}
		if indexUpdateNeeded {
			indexUpdates = append(indexUpdates, origPath)
		}
	}
	return indexUpdates, fileUpdates, nil
}

// updateTemplateFiles checks if the default.css and default.js files need to be
// updated in the output dir, and updates them if necessary.
func updateTemplateFiles() error {
	templateFiles := []string{"default.css", "default.js", "folder.svg"}
	for _, file := range templateFiles {
		outputFile := filepath.Join(config.Output, file)
		inputFile := filepath.Join("templates", config.Template, file)
		_, err := os.Stat(outputFile)
		if err != nil {
			if os.IsNotExist(err) {
				// File doesn't exist, so we need to copy it
				err := copyFile(inputFile, outputFile)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}

	return nil
}

// process is the main function that processes the original content, and generates
// the output.
func process() error {
	outputContent, err := getOutputContent()
	if err != nil {
		if strings.HasSuffix(err.(*os.PathError).Err.Error(), "no such file or directory") {
			log.Default().Printf("Creating output directory \"%s\"", config.Output)
			err = os.MkdirAll(config.Output, 0755)
			if err != nil {
				log.Fatalf("Failed to create output directory: %v", err)
			}
		}
	}
	originalContent, err := getOriginalContent()
	if err != nil {
		log.Fatalf("Failed to get original content: %v", err)
	}

	indexUpdates, fileUpdates, err := getUpdates(originalContent, outputContent)
	if err != nil {
		return err
	}

	err = processHTML(indexUpdates, &originalContent)
	if err != nil {
		return err
	}

	err = processImages(fileUpdates)
	if err != nil {
		return err
	}

	err = updateTemplateFiles()
	if err != nil {
		return err
	}

	return nil
}
