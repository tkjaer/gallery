package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddDir(t *testing.T) {
	dirMap := DirMap{}

	// Add a new directory
	dirMap.AddDir("/path/to/dir", "dir", true)

	// Verify the directory was added
	assert.Contains(t, dirMap, "/path/to/dir")
	assert.Equal(t, "dir", dirMap["/path/to/dir"].Name)
	assert.Equal(t, true, dirMap["/path/to/dir"].NeedsUpdate)

	// Add the same directory again
	dirMap.AddDir("/path/to/dir", "dir", false)

	// Verify the directory was not overwritten
	assert.Equal(t, true, dirMap["/path/to/dir"].NeedsUpdate)
}
