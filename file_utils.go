package main

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

// checkOrCreateOutputDir checks if the output directory exists, and creates it if it doesn't.
func checkOrCreateOutputDir() error {
	slog.Debug("Checking output directory", "output", config.Output)
	_, err := os.Stat(config.Output)
	if err != nil {
		if os.IsNotExist(err) {
			slog.Info("Output directory does not exist, creating it", "output", config.Output)
			err = os.MkdirAll(config.Output, 0755)
			if err != nil {
				return err
			}
			slog.Debug("Output directory created", "output", config.Output)
		} else {
			return err
		}
	}
	return nil
}

// copyFile copies a file from source to destination.
func copyFile(source, destination string) error {
	slog.Debug("Copying file", "source", source, "destination", destination)
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
	slog.Debug("File copied", "source", source, "destination", destination)
	return nil
}

// updateTemplateFiles checks if the default.css and default.js files need to be
// updated in the output dir, and updates them if necessary.
func updateTemplateFiles() error {
	slog.Debug("Updating template files")
	templateFiles := []string{"default.css", "default.js", "folder.svg"}
	for _, file := range templateFiles {
		slog.Debug("Processing template file", "file", file)
		outputFile := filepath.Join(config.Output, file)
		inputFile := filepath.Join("templates", config.Template, file)
		_, err := os.Stat(outputFile)
		if err != nil {
			if os.IsNotExist(err) {
				slog.Debug("Output file does not exist", "outputFile", outputFile)
				// File doesn't exist, so we need to copy it
				err := copyFile(inputFile, outputFile)
				if err != nil {
					return err
				}
				slog.Debug("Template file copied", "inputFile", inputFile, "outputFile", outputFile)
			} else {
				return err
			}
		}
	}

	return nil
}
