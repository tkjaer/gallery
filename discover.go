package main

import (
	"io/fs"
	"log/slog"
	"path/filepath"
	"strings"
	"time"
)

// File represents a file on disk, with a name and a modification time.
type File struct {
	Name    string
	ModTime time.Time
}

// SubDir represents a subdirectory on disk, with a name.
type SubDir struct {
	Name string
}

// DirContent represents the content of a directory on disk.
type Dir struct {
	Name    string
	ModTime time.Time
	Files   map[string]File
	SubDirs map[string]SubDir
}

// DirMap is a map of directories on disk, with the path as the key, and a list
// of files and subdirectories as the value.
type DirMap map[string]Dir

// AddDir adds a new directory to the DirMap. If the directory already exists, it does nothing.
func (dm DirMap) AddDir(path string, name string, modTime time.Time) {
	if _, exists := dm[path]; !exists {
		slog.Debug("Adding directory", "path", path, "name", name, "modTime", modTime)
		dm[path] = Dir{
			Name:    name,
			ModTime: modTime,
			Files:   make(map[string]File),
			SubDirs: make(map[string]SubDir),
		}
	}
}

// outputMap is a map of the output directory, with the path as the key, and the
// modification time as the value.
type outputMap map[string]time.Time

// getOutputContent returns a map of the output directory, with the path as the
// key, and the modification time as the value.
func getOutputContent() (outputMap, error) {
	slog.Debug("Getting output content", "output", config.Output)
	results := outputMap{}
	err := filepath.WalkDir(config.Output, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		fileInfo, err := d.Info()
		if err != nil {
			return err
		}
		results[path] = fileInfo.ModTime()
		return nil
	})
	return results, err
}

// getOriginalContent returns a map of the originals directory, with the path as the
// key, and a list of files and subdirectories as the value.
func getOriginalContent() (DirMap, error) {
	slog.Debug("Getting original content", "originals", config.Originals)
	results := DirMap{}
	err := filepath.WalkDir(config.Originals, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		slog.Debug("Processing file", "path", path, "name", d.Name())
		name := d.Name()
		parentDir := filepath.Dir(path)
		fileInfo, err := d.Info()
		if err != nil {
			return err
		}
		modTime := fileInfo.ModTime()

		if d.IsDir() {
			results.AddDir(path, name, modTime)
			// As filepath.WalkDir returns files in lexicographical order, we can
			// safely assume that the parent directory has already been added to
			// the results map, unless the parent directory is config.Originals.
			if _, ok := results[parentDir]; ok {
				slog.Debug("Adding subdirectory", "path", path, "name", name)
				results[parentDir].SubDirs[path] = SubDir{
					Name: name,
				}
			}
		} else {
			if strings.HasSuffix(name, ".jpg") || strings.HasSuffix(name, ".jpeg") {
				slog.Debug("Adding file", "path", path, "name", name)
				results[parentDir].Files[path] = File{
					Name:    name,
					ModTime: modTime,
				}
			}
		}
		return nil
	})
	return results, err
}
