package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckOrCreateOutputDir(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	config.Output = filepath.Join(tempDir, "output")

	// Ensure the directory does not exist initially
	_, err := os.Stat(config.Output)
	assert.True(t, os.IsNotExist(err))

	// Call the function
	err = checkOrCreateOutputDir()
	assert.NoError(t, err)

	// Verify the directory was created
	_, err = os.Stat(config.Output)
	assert.NoError(t, err)
	assert.True(t, isDir(config.Output))
}

func TestCopyFile(t *testing.T) {
	// Create a temporary source and destination file
	srcFile, err := ioutil.TempFile("", "src")
	assert.NoError(t, err)
	defer os.Remove(srcFile.Name())

	destFile, err := ioutil.TempFile("", "dest")
	assert.NoError(t, err)
	defer os.Remove(destFile.Name())

	// Write some content to the source file
	content := []byte("test content")
	_, err = srcFile.Write(content)
	assert.NoError(t, err)
	srcFile.Close()

	// Call the function
	err = copyFile(srcFile.Name(), destFile.Name())
	assert.NoError(t, err)

	// Verify the content of the destination file
	destContent, err := ioutil.ReadFile(destFile.Name())
	assert.NoError(t, err)
	assert.Equal(t, content, destContent)
}

func TestUpdateTemplateFiles(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	config.Output = filepath.Join(tempDir, "output")
	config.Template = "default"

	// Create the output directory
	err := os.MkdirAll(config.Output, 0755)
	assert.NoError(t, err)

	// Create a temporary templates directory
	templatesDir := filepath.Join(tempDir, "templates", config.Template)
	err = os.MkdirAll(templatesDir, 0755)
	assert.NoError(t, err)

	// Create template files with unique content
	templateFiles := map[string]string{
		"default.css": ":root {",
		"default.js":  "document.addEventListener('keydown', function(event) {",
		"folder.svg":  "<?xml version=\"1.0\" encoding=\"UTF-8\" standalone=\"no\"?>",
	}
	for file, content := range templateFiles {
		err := os.WriteFile(filepath.Join(templatesDir, file), []byte(content), 0644)
		assert.NoError(t, err)
	}

	// Call the function
	err = updateTemplateFiles()
	assert.NoError(t, err)

	// Verify the files were copied to the output directory and contain the correct content
	for file, originalContent := range templateFiles {
		outputFile := filepath.Join(config.Output, file)
		_, err := os.Stat(outputFile)
		assert.NoError(t, err)

		// Verify the content matches the original
		copiedContent, err := os.ReadFile(outputFile)
		assert.NoError(t, err)
		assert.Contains(t, string(copiedContent), originalContent)
	}
}

// Helper function to check if a path is a directory
func isDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
