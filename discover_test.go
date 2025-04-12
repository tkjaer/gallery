package main

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTempDirStructure(t *testing.T, baseDir string, structure []map[string]any) {
	for _, item := range structure {
		name, ok := item["name"].(string)
		assert.True(t, ok, "name must be a string")
		path := filepath.Join(baseDir, name)

		if content, isDir := item["content"].([]map[string]any); isDir {
			log.Default().Println("Creating directory:", path)
			err := os.Mkdir(path, 0755)
			assert.NoError(t, err)
			createTempDirStructure(t, path, content) // Recursively create subdirectory structure
		} else if content, isFile := item["content"].(string); isFile {
			log.Default().Println("Creating file:     ", path)
			err := os.WriteFile(path, []byte(content), 0644)
			assert.NoError(t, err)
		} else {
			t.Fatalf("Invalid content type for %s", name)
		}
	}
}

func TestGetOutputContent(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	config.Output = tempDir

	// Create a mock directory structure
	structure := []map[string]any{
		{
			"name":    "file1.txt",
			"content": "content1",
		},
		{
			"name":    "file2.jpg",
			"content": "content2",
		},
		{
			"name":    "subdir",
			"content": []map[string]any{},
		},
		{
			"name":    "subdir/file3.jpg",
			"content": "content3",
		},
	}
	createTempDirStructure(t, tempDir, structure)

	// Call getOutputContent
	output, err := getOutputContent()
	assert.NoError(t, err)

	// Verify the output
	assert.Contains(t, output, filepath.Join(tempDir, "file1.txt"))
	assert.Contains(t, output, filepath.Join(tempDir, "file2.jpg"))
	assert.Contains(t, output, filepath.Join(tempDir, "subdir"))
	assert.Contains(t, output, filepath.Join(tempDir, "subdir/file3.jpg"))
}

func TestGetOriginalContent(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	config.Originals = tempDir

	// Create a mock directory structure
	structure := []map[string]any{
		{
			"name":    "file1.txt",
			"content": "content1",
		},
		{
			"name":    "file2.jpg",
			"content": "content2",
		},
		{
			"name":    "subdir",
			"content": []map[string]any{},
		},
		{
			"name":    "subdir/file3.jpg",
			"content": "content3",
		},
	}
	createTempDirStructure(t, tempDir, structure)

	// Call getOriginalContent
	originals, err := getOriginalContent()
	assert.NoError(t, err)

	// Verify the output
	assert.Contains(t, originals, tempDir)
	assert.Contains(t, originals[filepath.Join(tempDir, "subdir")].Files, filepath.Join(tempDir, "subdir/file3.jpg"))
	assert.NotContains(t, originals[filepath.Join(tempDir, "subdir")].Files, filepath.Join(tempDir, "file1.txt"))
}
